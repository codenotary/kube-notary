/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package watcher

import (
	"fmt"
	"time"

	"github.com/vchain-us/kube-notary/pkg/config"
	"github.com/vchain-us/kube-notary/pkg/metrics"
	"github.com/vchain-us/kube-notary/pkg/verify"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	Run()
}

type watchdog struct {
	clientset *kubernetes.Clientset
	log       *log.Logger
	rec       metrics.Recorder
	cfg       config.Interface
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
	}, nil
}

func (w *watchdog) Run() {
	clientset := w.clientset
	for {
		w.log.SetLevel(w.cfg.LogLevel())

		ns := w.cfg.Namespace()
		sleep := w.cfg.Interval()
		keys := w.cfg.TrustedKeys()
		fields := log.Fields{
			"namespace":   ns,
			"interval":    sleep,
			"trustedKeys": keys,
		}

		pods, err := clientset.CoreV1().Pods(ns).List(metav1.ListOptions{})
		if err != nil {
			fields["error"] = true
			w.log.WithFields(fields).Errorf("Error getting pods: %s", err)
		} else {
			fields["podCount"] = len(pods.Items)
			w.log.WithFields(fields).Debug("Verification started")

			for _, pod := range pods.Items {
				w.watchPod(pod, keys...)
			}
		}

		w.log.Debugf("Sleeping for %s", sleep)
		time.Sleep(sleep)
	}
}

func (w *watchdog) watchPod(pod corev1.Pod, trustedKeys ...string) {
	for _, status := range pod.Status.ContainerStatuses {
		if status.State.Running == nil {
			w.log.Infof(`Container "%s" in pod "%s" is not running: skipped`, status.Name, pod.Name)
			continue
		}

		hash, verification, err := verify.ImageID(status.ImageID, trustedKeys...)
		if err != nil {
			w.log.Errorf(`Cannot verify "%s" in pod "%s": %s`, status.ImageID, pod.Name, err)
			continue
		}

		metric := metrics.Metric{
			Pod:             &pod,
			ContainerStatus: &status,
			Verification:    verification,
			Hash:            hash,
		}

		fields := metric.LogFields()
		if verification.Trusted() {
			w.log.WithFields(*fields).Info("Image is trusted")
		} else {
			w.log.WithFields(*fields).Warn("Image is NOT trusted")
		}

		w.rec.Record(metric)
	}
}
