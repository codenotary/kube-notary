/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/vchain-us/kube-notary/pkg/config"
	"github.com/vchain-us/kube-notary/pkg/metrics"
	"github.com/vchain-us/kube-notary/pkg/watcher"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	// creates the in-cluster config
	clusterCfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(clusterCfg)
	if err != nil {
		panic(err.Error())
	}
	// creates the logger
	logger := logrus.New()
	// creates the metrics recorder
	recorder := metrics.NewRecorder()
	// creates the watcher configuration
	cfg, err := config.New()
	if err != nil {
		panic(err.Error())
	}
	// creates and run the watcher
	w, err := watcher.New(clientset, cfg, recorder, logger)
	if err != nil {
		panic(err.Error())
	}
	go w.Run()

	// The metrics.Handler provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", metrics.Handler())

	// The w.ResultsHandler provides a handler to expose detailed
	// verification results.
	http.Handle("/results", w.ResultsHandler())

	// Healthcheck endpoint.
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	panic(http.ListenAndServe(":9581", nil))
}
