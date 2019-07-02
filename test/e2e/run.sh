#!/usr/bin/env sh

set -euo pipefail

CLUSTER_NAME=kube-notary-test

# Clean up an existing cluster instance eventually
echo $(kind delete cluster --name=$CLUSTER_NAME 2>&1) > /dev/null

# Create the cluster
kind create cluster --name=$CLUSTER_NAME
kind load docker-image --name=$CLUSTER_NAME kube-notary:test
KUBECONFIG=$(kind get kubeconfig-path --name=$CLUSTER_NAME)
CONTEXT=kubernetes-admin@$CLUSTER_NAME

# Setup tiller service account and init Helm
kubectl --kubeconfig $KUBECONFIG --context $CONTEXT apply -f ./tiller-rbac.yaml
helm init --kubeconfig $KUBECONFIG --kube-context $CONTEXT --service-account tiller --history-max 200 --wait

# Install kube-notary chart
helm install --kubeconfig $KUBECONFIG --kube-context $CONTEXT \
    -n kube-notary ../../helm/kube-notary \
    --set image.repository=kube-notary --set image.tag=test \
    --set image.pullPolicy=Never \
    --wait


# Tests
sleep 1
echo "Running tests..."
kubectl --kubeconfig $KUBECONFIG --context $CONTEXT get pods --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary"
kubectl --kubeconfig $KUBECONFIG --context $CONTEXT get service --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary"
echo "Tests done."

# Teardown the cluster
kind delete cluster --name=$CLUSTER_NAME