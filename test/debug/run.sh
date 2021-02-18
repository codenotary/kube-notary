#!/bin/bash

# Copyright (c) 2019 vChain, Inc. All Rights Reserved.
# This software is released under GPL3.
# The full license information can be found under:
# https://www.gnu.org/licenses/gpl-3.0.en.html

# Creates a local kube-notary cluster with kind for debug pupose.
# A debug image with a delve server is launched. Debugger is bind at 40000 port.

set -euo pipefail

CLUSTER_NAME=kube-notary-dbg
KUBE_NOTARY_IMAGE=kube-notary:debug

# Clean up an existing cluster instance eventually
echo $(kind delete cluster --name=$CLUSTER_NAME 2>&1) > /dev/null

# Create the cluster
kind create cluster --name=$CLUSTER_NAME
kubectl cluster-info --context kind-$CLUSTER_NAME
kind load docker-image --name=$CLUSTER_NAME $KUBE_NOTARY_IMAGE

# Setup tiller service account and init Helm
kubectl  apply -f ../../test/e2e/tiller-rbac.yaml
helm init --service-account tiller --history-max 200 --wait

# Not needed in CodeNotary.io mode
kubectl create secret generic vcn-lc-api-key --from-literal=api-key=trqgnxwyjdwmcuajmczcrtjccagzhiawzkod

# Install kube-notary chart
# Remove cnlc.host to disable Ledger Compliance mode
helm install \
    -n kube-notary ../../helm/kube-notary \
    --set debug=true \
    --set image.repository=kube-notary --set image.tag=debug\
    --set image.pullPolicy=Never \
    --set restartPolicy=Never \
    --set cnlc.host=$(hostname) \
    --set watch.interval=10s \
    --wait

export SERVICE_NAME=service/$(kubectl get service --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward --namespace default $SERVICE_NAME 9581 &
kubectl port-forward --namespace default $SERVICE_NAME 40000 &
