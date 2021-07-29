#!/bin/bash

# Copyright (c) 2019 vChain, Inc. All Rights Reserved.
# This software is released under GPL3.
# The full license information can be found under:
# https://www.gnu.org/licenses/gpl-3.0.en.html

# Creates a local kube-notary cluster with kind for debug pupose.
# A debug image with a delve server is launched. Debugger is bind at 40000 port.
# Please make sure helm 3 is installed

set -euo pipefail

CLUSTER_NAME=kube-notary-dbg
KUBE_NOTARY_IMAGE=kube-notary:debug

# Clean up an existing cluster instance eventually
echo $(kind delete cluster --name=$CLUSTER_NAME 2>&1) > /dev/null

# Create the cluster
kind create cluster --name=$CLUSTER_NAME
kubectl cluster-info --context kind-$CLUSTER_NAME
kind load docker-image --name=$CLUSTER_NAME $KUBE_NOTARY_IMAGE

# HELM v2 Setup tiller service account and init Helm
# kubectl  apply -f ../../test/e2e/tiller-rbac.yaml
# helm init --service-account tiller --history-max 200 --wait

# Not needed in CodeNotary.io mode
kubectl create secret generic vcn-lc-api-key --from-literal=api-key=kube.izZQyDUfWnOSwSZefrOUThcVdbGBOouzKnHf

# Install kube-notary chart
# Remove cnlc.host to disable Ledger Compliance mode
# When debugging CNLC on the same network(hostNetwork: true in deployment config) cnlc.host need to be set to the docker bridge ip.
# To retrieve it is possible to use:
# docker network inspect bridge -f '{{range .IPAM.Config}}{{.Gateway}}{{end}}'
# or
# /sbin/ip route|awk 'FNR == 6 {print $9}'
helm install \
    -n default kubeinstance ../../helm/kube-notary \
    --set debug=true \
    --set image.repository=kube-notary --set image.tag=debug\
    --set image.pullPolicy=Never \
    --set restartPolicy=Never \
    --set cnlc.host=$(docker network inspect bridge -f '{{range .IPAM.Config}}{{.Gateway}}{{end}}') \
    --set watch.interval=10s \
    --wait

export SERVICE_NAME=service/kubeinstance-kube-notary
kubectl port-forward $SERVICE_NAME 9581 &
kubectl port-forward $SERVICE_NAME 40000 &
