/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package watcher

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/codenotary/vcn-enterprise/pkg/api"
	"github.com/codenotary/vcn-enterprise/pkg/meta"
	"github.com/vchain-us/kube-notary/pkg/config"
	"github.com/vchain-us/kube-notary/pkg/image"
	"github.com/vchain-us/kube-notary/pkg/metrics"
	"github.com/vchain-us/kube-notary/pkg/verify"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const kubeNotaryWatcherName = "kube-notary"
const kubeSystemNamespace = "kube-system"

type WatchDog struct {
	clientSet  *kubernetes.Clientset
	rec        metrics.Recorder
	cfg        *config.Config
	res        map[string]Result
	tmp        []string
	idx        []string
	seen       map[string]bool
	imageCache map[string]string
	mu         *sync.RWMutex
}

func New(clientSet *kubernetes.Clientset, cfg *config.Config, rec metrics.Recorder) *WatchDog {
	return &WatchDog{
		clientSet:  clientSet,
		rec:        rec,
		cfg:        cfg,
		res:        map[string]Result{},
		tmp:        []string{},
		idx:        []string{},
		seen:       map[string]bool{},
		imageCache: map[string]string{},
		mu:         &sync.RWMutex{},
	}
}

func (w *WatchDog) Run() {
	log.Infof("WatchDog started on namespace %s interval %s LcHost %s Port %s LedgerName %s", w.cfg.Namespace(), w.cfg.Interval(), w.cfg.LcHost(), w.cfg.LcPort(), w.cfg.LcCrossLedgerKeyLedgerName())

	keys := w.cfg.TrustedKeys()
	org := w.cfg.TrustedOrg()

	var opt verify.Option
	if org != "" {
		opt = verify.WithSignerOrg(org)
		if len(keys) > 0 {
			log.Warn("Trusted keys ignored because an organization is set")
			keys = nil
		}
	} else if len(keys) > 0 {
		opt = verify.WithSignerKeys(keys...)
	}

	for {
		pods, err := w.clientSet.CoreV1().Pods(w.cfg.Namespace()).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			log.Errorf("Error getting pods: %s", err)
			continue
		}

		statusMap := make(map[meta.Status]int, 5)
		for _, pod := range pods.Items {
			statuses := w.watchPod(pod, opt)
			for i, s := range statuses {
				statusMap[i] += s
			}
		}
		for st, c := range statusMap {
			metrics.SetTotals(w.cfg.Namespace(), st.String(), c)
		}
		w.commit()
		time.Sleep(w.cfg.Interval())
	}
}

func (w *WatchDog) watchPod(pod corev1.Pod, options ...verify.Option) (statuses map[meta.Status]int) {
	statuses = make(map[meta.Status]int, 5)
	log.Infof("Processing Pod %s:%s", pod.Namespace, pod.Name)

	// skip K8s watcher container
	if strings.Contains(pod.Name, kubeNotaryWatcherName) || strings.Contains(pod.Namespace, kubeSystemNamespace) {
		return
	}

	pullSecrets := make([]string, len(pod.Spec.ImagePullSecrets))
	for i, localRef := range pod.Spec.ImagePullSecrets {
		pullSecrets[i] = localRef.Name
	}

	keychain, err := image.NewKeychain(
		w.clientSet,
		pod.Namespace,
		pod.Spec.ServiceAccountName,
		pullSecrets,
	)
	if err != nil {
		log.Warnf(`Keychain error in pod "%s": %s`, pod.Name, err)
	}

	// make options
	l := len(options) + 1
	opts := make([]verify.Option, len(options)+1)
	copy(opts, options)
	opts[l-1] = verify.WithAuthKeychain(keychain)
	for _, status := range pod.Status.ContainerStatuses {
		v := &verify.Verification{}

		if status.State.Running == nil {
			log.Infof(`Container "%s" in pod "%s" is not running: skipped`, status.Name, pod.Name)
		}
		errorList := make([]error, 0)

		if status.ImageID == "" {
			continue
		}

		var hash string
		var err error
		var ok bool
		if hash, ok = w.getAuthorized(status.ImageID); !ok {
			log.Infof("getting Image Hash from %s ", status.Name)
			hash, err = verify.ImageHash(
				status.ImageID,
				opts...,
			)
			if err != nil {
				errorList = append(errorList, err)
				v.Status = meta.StatusUnknown
				v.Level = meta.LevelUnknown
				v.Date = ""
				v.Trusted = false
				errorList = append(errorList, err)
				log.Errorf(`Unable to get ImageHash from Registry "%s" in pod "%s": %s`, status.ImageID, pod.Name, err)
			}

			if hash != "" && err == nil {
				w.setAuthorized(status.ImageID, hash)
			}
		}

		if hash == "" {
			log.Errorf("Cannot Veryfy from empty HASH image name %s id %s", status.Image, status.ImageID)
			continue
		}

		log.Infof("Veryfy image name %s id %s hash %s", status.Image, status.ImageID, hash)

		apiKey, apiKeyErr := w.cfg.ApiKey() // @TODO: To init App
		if apiKeyErr != nil {
			log.Warnf("Unable to get Api Key from config, error: %v", apiKeyErr)
			return
		}
		ar, err := VerifyArtifact(hash, apiKey, w.cfg.LcCrossLedgerKeyLedgerName(), w.cfg.LcSignerID(), w.cfg.LcHost(), w.cfg.LcPort(), w.cfg.LcCert(), w.cfg.LcSkipTlsVerify(), w.cfg.LcNoTls())

		if errors.Is(err, api.ErrNotVerified) {
			v.Status = meta.StatusUnknown
			v.Level = meta.LevelUnknown
			v.Date = ""
			v.Trusted = false
			log.Errorf("Image %s in pod %s is not verified: %s", status.ImageID, pod.Name, err)
		} else if errors.Is(err, api.ErrNotFound) {
			v.Status = meta.StatusUnknown
			v.Level = meta.LevelUnknown
			v.Date = ""
			v.Trusted = false
			log.Errorf("Image %s in pod %s not found: %s", status.ImageID, pod.Name, err)
		} else if err == nil {
			v.Status = ar.Status
			v.Level = meta.LevelCNLC
			v.Date = ar.Date()
			v.Trusted = false
			if ar.Status == meta.StatusTrusted {
				v.Trusted = true
			}
			log.Infof("Image %s with ID %s is trusted", status.Image, status.ImageID)
		} else {
			v.Status = meta.StatusUnknown
			v.Level = meta.LevelUnknown
			v.Date = ""
			v.Trusted = false
			errorList = append(errorList, err)
			log.Errorf("Cannot verify %s in pod %s: %s", status.ImageID, pod.Name, err)
		}
		statuses[v.Status]++
		// @TODO: Record metric for a pod ¿?¿?
		w.rec.Record(metrics.Metric{
			Pod:             &pod,
			ContainerStatus: &status,
			Verification:    v,
			Hash:            hash,
		})

		// update or insert the result into tmp list
		w.upsert(pod, status, v, hash, errorList)
	}

	return
}

func VerifyArtifact(hash, apiKey, lcLedger, signerID, lcHost, lcPort, lcCert string, lcSkipTlsVerify, lcNoTls bool) (a *api.LcArtifact, err error) {

	log.Printf("VerifyArtifact apiKey %s ledger %s host %s port %s cert %s skip %v noTls %v \n", apiKey, lcLedger, lcHost, lcPort, lcCert, lcSkipTlsVerify, lcNoTls)

	cl, err := buildClient(apiKey, lcLedger, lcHost, lcPort, lcCert, lcSkipTlsVerify, lcNoTls)
	if err != nil {
		return nil, fmt.Errorf("unable to build client, error %w", err)
	}

	hash = strings.TrimPrefix(hash, "sha256:")
	metadata := map[string][]string{meta.VcnLCCmdHeaderName: {meta.VcnLCVerifyCmdHeaderValue}}
	a, _, err = cl.LoadArtifact(hash, signerID, "", 0, metadata)
	if errors.Is(err, api.ErrNotFound) {
		return nil, fmt.Errorf("no artifact found on hash %s, error %w", hash, err)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to load artifact on hash %s, error %w", hash, err)
	}

	return a, err
}

func buildClient(apiKey, lcLedger, lcHost, lcPort, lcCert string, lcSkipTlsVerify, lcNoTls bool) (*api.LcUser, error) {
	client, err := api.NewLcClient(apiKey, lcLedger, lcHost, lcPort, lcCert, lcSkipTlsVerify, lcNoTls, nil, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("unable to create DataService client, error %w", err)
	}
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("unable to connect dataService, error %w", err)
	}
	return &api.LcUser{
		Client: client,
	}, nil
}
