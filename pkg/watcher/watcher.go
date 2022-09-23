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
	"fmt"
	"net/http"
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

type Interface interface {
	Run()
	ResultsHandler() http.Handler
}

type watchdog struct {
	clientset *kubernetes.Clientset
	log       *log.Logger
	rec       metrics.Recorder
	cfg       config.Interface
	res       map[string]Result
	tmp       []string
	idx       []string
	seen      map[string]bool
	mu        *sync.RWMutex
}

func New(clientset *kubernetes.Clientset, cfg config.Interface, rec metrics.Recorder, logger *log.Logger) (Interface, error) {

	if clientset == nil {
		return nil, fmt.Errorf("clientset cannot be nil")
	}

	if logger == nil {
		logger = log.StandardLogger()
	}

	return &watchdog{
		clientset: clientset,
		log:       logger,
		rec:       rec,
		cfg:       cfg,
		res:       map[string]Result{},
		tmp:       []string{},
		idx:       []string{},
		seen:      map[string]bool{},
		mu:        &sync.RWMutex{},
	}, nil
}

func (w *watchdog) Run() {
	clientset := w.clientset
	for {
		w.log.SetLevel(w.cfg.LogLevel())

		ns := w.cfg.Namespace()
		sleep := w.cfg.Interval()
		keys := w.cfg.TrustedKeys()
		org := w.cfg.TrustedOrg()
		fields := log.Fields{
			config.LogLevel:       w.cfg.LogLevel().String(),
			config.WatchNamespace: ns,
			config.WatchInterval:  sleep,
			config.TrustKeys:      keys,
		}

		var opt verify.Option

		if org != "" {
			delete(fields, config.TrustKeys)
			fields[config.TrustOrg] = org
			opt = verify.WithSignerOrg(org)
			if len(keys) > 0 {
				w.log.WithFields(fields).Warn("Trusted keys ignored because an organization is set")
				keys = nil
			}
		} else if len(keys) > 0 {
			opt = verify.WithSignerKeys(keys...)
		}

		w.rec.Reset()

		pods, err := clientset.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			fields["error"] = true
			w.log.WithFields(fields).Errorf("Error getting pods: %s", err)
		} else {
			fields["podCount"] = len(pods.Items)
			w.log.WithFields(fields).Debug("Verification started")

			for _, pod := range pods.Items {
				w.watchPod(pod, opt)
			}
		}

		// commit tmp list into results index
		w.commit()
		w.log.Debugf("Sleeping for %s", sleep)
		time.Sleep(sleep)
	}
}

func (w *watchdog) watchPod(pod corev1.Pod, options ...verify.Option) {

	pullSecrets := make([]string, len(pod.Spec.ImagePullSecrets))
	for i, localRef := range pod.Spec.ImagePullSecrets {
		pullSecrets[i] = localRef.Name
	}

	keychain, err := image.NewKeychain(
		w.clientset,
		pod.Namespace,
		pod.Spec.ServiceAccountName,
		pullSecrets,
	)
	if err != nil {
		w.log.Warnf(`Keychain error in pod "%s": %s`, pod.Name, err)
	}

	// make options
	l := len(options) + 1
	opts := make([]verify.Option, len(options)+1)
	copy(opts, options)
	opts[l-1] = verify.WithAuthKeychain(keychain)

	for _, status := range pod.Status.ContainerStatuses {
		v := &verify.Verification{}

		if status.State.Running == nil {
			w.log.Infof(`Container "%s" in pod "%s" is not running: skipped`, status.Name, pod.Name)
		}
		errorList := make([]error, 0)

		hash, err := verify.ImageHash(
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
			w.log.Errorf(`Cannot verify "%s" in pod "%s": %s`, status.ImageID, pod.Name, err)
		}

		if w.cfg.LcHost() != "" && hash != "" {
			hash = strings.TrimPrefix(hash, "sha256:")
			// @TODO: Verify hash against CNCL, returns CNLC artifact (SBOM)
			ar, err := api.PublicCNLCVerify(hash, w.cfg.LcCrossLedgerKeyLedgerName(), w.cfg.LcSignerID(), w.cfg.LcHost(), w.cfg.LcPort(), w.cfg.LcCert(), w.cfg.LcSkipTlsVerify(), w.cfg.LcNoTls())
			metric := metrics.Metric{
				Pod:             &pod,
				ContainerStatus: &status,
				Verification:    v,
				Hash:            hash,
			}
			fields := metric.LogFields()
			switch err {
			case api.ErrNotVerified: // @TODO: Not verified
				v.Status = meta.StatusUnknown
				v.Level = meta.LevelUnknown
				v.Date = ""
				v.Trusted = false
				w.log.Warnf(`Image "%s" in pod "%s" is not verified: %s`, status.ImageID, pod.Name, err)
			case api.ErrNotFound: // @TODO: Not found
				v.Status = meta.StatusUnknown
				v.Level = meta.LevelUnknown
				v.Date = ""
				v.Trusted = false
				w.log.Warnf(`Image "%s" in pod "%s" not found: %s`, status.ImageID, pod.Name, err)
			case nil:
				v.Status = ar.Status
				v.Level = meta.LevelCNLC
				v.Date = ar.Date()
				v.Trusted = false
				if ar.Status == meta.StatusTrusted {
					v.Trusted = true
				}
				w.log.WithFields(*fields).Info("Image is trusted") // @TODO: On ar.Status different from trusted it will be false trust positive
			default:
				v.Status = meta.StatusUnknown
				v.Level = meta.LevelUnknown
				v.Date = ""
				v.Trusted = false
				errorList = append(errorList, err)
				w.log.Errorf(`Cannot verify "%s" in pod "%s": %s`, status.ImageID, pod.Name, err)
			}
			w.rec.Record(metric)

		}
		// update or insert the result into tmp list
		w.upsert(pod, status, v, hash, errorList)
	}
}
