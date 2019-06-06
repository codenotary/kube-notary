/*
 * Copyright (c) 2019 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 *
 */

package main

import (
	"github.com/sirupsen/logrus"

	"github.com/vchain-us/kubewatch/pkg/config"
	"github.com/vchain-us/kubewatch/pkg/watcher"

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
	// creates the watcher configuration
	cfg, err := config.New()
	if err != nil {
		panic(err.Error())
	}
	// creates and run the watcher
	w, err := watcher.New(clientset, cfg, logger)
	if err != nil {
		panic(err.Error())
	}
	w.Run()
}
