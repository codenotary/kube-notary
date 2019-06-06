/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package watcher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vchain-us/kubewatch/pkg/verify"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	Run() error
}

type watchdog struct {
	clientset *kubernetes.Clientset
	log       *log.Logger
	cfg       *Config
}

type Config struct {
	Namespace string
	Interval  time.Duration
}

func New(clientset *kubernetes.Clientset, config Config, logger *log.Logger) (Interface, error) {

	if clientset == nil {
		return nil, fmt.Errorf("clientset cannot be nil")
	}

	if logger == nil {
		logger = log.StandardLogger()
	}

	return &watchdog{
		clientset: clientset,
		log:       logger,
		cfg:       &config,
	}, nil
}

func (w *watchdog) Run() error {
	clientset := w.clientset
	for {
		pods, err := clientset.CoreV1().Pods(w.cfg.Namespace).List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		w.log.Infof("There are %d pods in the cluster", len(pods.Items))

		for _, pod := range pods.Items {
			w.watchPod(pod)
		}

		time.Sleep(w.cfg.Interval)
	}
}

func (w *watchdog) watchPod(pod corev1.Pod) {
	for _, status := range pod.Status.ContainerStatuses {

		verification, err := verify.ImageID(status.ImageID)
		if err != nil {
			w.log.Errorf("Cannot verify %s in pod %s: %s", status.ImageID, pod.Name, err)
			continue
		}

		b, _ := json.Marshal(verification)

		fields := log.Fields{
			"pod":          pod.Name,
			"image":        status.Image,
			"imageID":      status.ImageID,
			"verification": string(b),
			"status":       verification.Status,
			"trusted":      verification.Trusted(),
		}

		if verification.Trusted() {
			w.log.WithFields(fields).Infof("Image %s (digest: %s) is trusted", status.Image, status.ImageID)
		} else {
			w.log.WithFields(fields).Warnf("Image %s (digest: %s) is NOT trusted", status.Image, status.ImageID)
		}
	}
}
