#!/bin/bash

# Copyright (c) 2019 vChain, Inc. All Rights Reserved.
# This software is released under GPL3.
# The full license information can be found under:
# https://www.gnu.org/licenses/gpl-3.0.en.html

# Creates a local kube-notary cluster with kind.
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
kubectl  apply -f ../../test/e2e/tiller-rbac.yaml
helm init --service-account tiller --history-max 200 --wait

# Install kube-notary chart
helm install \
    -n kube-notary ../../helm/kube-notary \
    --set image.repository=$KUBE_NOTARY_IMAGE --set image.tag=$KUBE_NOTARY_TAG \
    --set image.pullPolicy=Never \
    --wait

export SERVICE_NAME=service/$(kubectl get service --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward --namespace default $SERVICE_NAME 9581