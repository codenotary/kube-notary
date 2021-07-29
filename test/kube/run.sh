#!/bin/bash

# Copyright (c) 2019 vChain, Inc. All Rights Reserved.
# This software is released under GPL3.
# The full license information can be found under:
# https://www.gnu.org/licenses/gpl-3.0.en.html

# Creates a local kube-notary cluster with kind.
# Please make sure helm 3 is installed

set -euo pipefail

CLUSTER_NAME=kube-notary-cluster
KUBE_NOTARY_IMAGE=codenotary/kube-notary
KUBE_NOTARY_TAG=latest

# Clean up an existing cluster instance eventually
echo $(kind delete cluster --name=$CLUSTER_NAME 2>&1) > /dev/null

# Create the cluster
kind create cluster --name=$CLUSTER_NAME
kubectl cluster-info --context kind-$CLUSTER_NAME
kind load docker-image --name=$CLUSTER_NAME $KUBE_NOTARY_IMAGE:$KUBE_NOTARY_TAG

# Setup tiller service account and init Helm
# kubectl  apply -f ../../test/e2e/tiller-rbac.yaml
# helm init --service-account tiller --history-max 200 --wait

# Not needed in CodeNotary.io mode
# kubectl create secret generic vcn-lc-api-key --from-literal=api-key=trqgnxwyjdwmcuajmczcrtjccagzhiawzkod

# Install kube-notary chart
helm install \
    -n default kubeinstance ../../helm/kube-notary \
    --set image.repository=$KUBE_NOTARY_IMAGE --set image.tag=$KUBE_NOTARY_TAG \
    --set image.pullPolicy=Never \
    --wait

export SERVICE_NAME=service/kubeinstance-kube-notary
kubectl port-forward --namespace default $SERVICE_NAME 9581
