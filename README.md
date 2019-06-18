# kube-notary
> A Kubernetes watchdog for verifying image trust with CodeNotary.

*Work in progress!*

## Install
`kube-notary` is installed using a Helm chart.
> Kubernetes 1.9 or above, and Helm 2.8 or above need to be installed in your cluster.

First, make sure your local config is pointing to the context you want to use (ie. check `kubectl config current-context`).
Then, to install `kube-notary`:

* Clone this repository locally.
* Change directory into `kube-notary`.
* Finally run:
```
helm install -n kube-notary helm/kube-notary
```

Alternatively, it's possible to manually install `kube-notary` without using Helm. Templates for manual installation are within the [kubernetes folder](kubernetes).

## Usage

`kube-notary` provides both detailed log output and a Prometheus metrics endpoint to monitor the verification status of your running containers. After the installation you will find instruction to get them.

Examples:
```
  # Metrics endpoint
  export SERVICE_NAME=service/$(kubectl get service --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
  echo "Check metrics endpoint at http://127.0.0.1:9581/metrics"
  kubectl port-forward --namespace default $SERVICE_NAME 9581

  # Stream logs
  export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
  kubectl logs --namespace default -f $POD_NAME
```
### Metrics

If a Prometheus installation is running within your cluster, metrics provided by `kube-notary` will be automatically discovered. 
Furthermore, you can find an example of preconfigured Grafana dashboard [here](grafana/).

## Configuration

By default, `kube-notary` is installed into the current namespace (you can change it by using `helm install --namespace`) but it will watch to pods in all namespaces.

At install time you can change any values of [helm/kube-notary/values.yaml](helm/kube-notary/values.yaml) by using the Helm's `--set` option.
For example, to instruct `kube-notary` to check only the `kube-system` namespace, just use:
```
helm install -n kube-notary helm/kube-notary --set watch.namespace="kube-system"
```

### Runtime configuration

The following options with [helm/kube-notary/values.yaml](helm/kube-notary/values.yaml) have effect on the `kube-notary` runtime behavior.
```
# Runtime config
log:
  level: debug
watch: 
  namespace: ""
  interval: 60s
```

During the installation, they are stored in a `configmap`. Configuration hot-reloading is supported, so you can modify and apply the configmap while `kube-notary` is running. 

For example, to change the watching interval from default to `30s`:
```
kubectl patch configmaps/kube-notary \
    --type merge \
    -p '{"data":{"config.yaml":"log:\n  level: debug\nwatch: \n  namespace: \n  interval: 30s"}}'
```

