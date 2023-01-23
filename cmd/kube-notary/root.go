/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vchain-us/kube-notary/pkg/config"
	"github.com/vchain-us/kube-notary/pkg/metrics"
	"github.com/vchain-us/kube-notary/pkg/status"
	"github.com/vchain-us/kube-notary/pkg/watcher"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

const httpPort = 9581

func main() {
	configFilePath := flag.String("config", config.DefaultConfigPath, "config file path")
	mode := flag.String("k8s-mode", config.InternalMode, "kubernetes controller mode, external runs out of the cluster (development mode")
	flag.Parse()

	cfg, err := config.New(*configFilePath)
	if err != nil {
		log.Fatalf("unable to load config, error %v", err)
	}

	log.Infof("kube-notary watcher started on namespace %s with watch interval %s, listening http calls on port %d", cfg.Namespace(), cfg.Interval(), httpPort)

	if *mode == config.InternalMode {
		clusterCfg, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalf("unable to get cluster config from flags, error %v", err)
		}
		run(cfg, clusterCfg)
		return
	}

	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"
	clusterCfg, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("unable to get cluster config from flags, error %v", err)
	}

	run(cfg, clusterCfg)
}

func run(cfg *config.Config, clusterCfg *rest.Config) {
	clientSet, err := kubernetes.NewForConfig(clusterCfg)
	if err != nil {
		log.Fatalf("unable to create kubernetes client, error %v", err)
	}

	w := watcher.New(clientSet, cfg, metrics.NewRecorder())
	go w.Run()

	http.Handle("/metrics", metrics.Handler())
	http.Handle("/results", w.ResultsHandler())
	http.Handle("/bulk_sign", w.BulkHandler())
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	http.Handle("/", status.Handler())

	panic(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil))
}
