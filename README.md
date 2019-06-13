# kubewatch
> A Kubernetes watchdog for verifying image trust with CodeNotary.

*Work in progress!*

## Install
Kubewatch is installed using a Helm chart.
> Kubernetes 1.9 or above, and Helm 2.8 or above need to be installed in your cluster.

First, make sure your local config is pointing to the context you want to use (ie. check `kubectl config current-context`).
Then, to install Kubewatch:

* Clone this repository locally.
* Change directory into the `kubewatch` Git repository.
* Finally run:
```
helm install -n kubewatch helm/kubewatch
```

## Usage

`kubewatch` provides both detailed log output and a Prometheus metrics endpoint to monitor the verification status of your running containers. After the installation you will find instruction to get them.

Examples:
```
  # Metrics endpoint
  export SERVICE_NAME=service/$(kubectl get service --namespace default -l "app.kubernetes.io/name=kubewatch,app.kubernetes.io/instance=kubewatch" -o jsonpath="{.items[0].metadata.name}")
  echo "Check metrics endpoint at http://127.0.0.1:9581/metrics"
  kubectl port-forward --namespace default $SERVICE_NAME 9581

  # Stream logs
  export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=kubewatch,app.kubernetes.io/instance=kubewatch" -o jsonpath="{.items[0].metadata.name}")
  kubectl logs --namespace default -f $POD_NAME
```
### Metrics

If a Prometheus installation is running within your cluster, metrics provided by `kubewatch` will be automatically discovered. 
Furthermore, you can find an example of preconfigured Grafana dashboard [here](grafana/).

## Configuration

By default, `kubewatch` is installed into the current namespace (you can change it by using `helm install --namespace`) but it will watch to pods in all namespaces.

At install time you can change any values of [helm/kubewatch/values.yaml](helm/kubewatch/values.yaml) by using the Helm's `--set` option.
For example, to instruct `kubewatch` to check only the `kube-system` namespace, just use:
```
helm install -n kubewatch helm/kubewatch watch.namespace="kube-system"
```

### Runtime configuration

The following options with [helm/kubewatch/values.yaml](helm/kubewatch/values.yaml) have effect on the `kubewatch` runtime behavior.
```
# Runtime config
log:
  level: debug
watch: 
  namespace: ""
  interval: 60s
```

During the installation, they are stored in a `configmap`. Configuration hot-reloading is supported, so you can modify and apply the configmap while `kubewatch` is running. 

For example, to change the watching interval from default to `30s`:
```
kubectl patch configmaps/kubewatch \
    --type merge \
    -p '{"data":{"config.yaml":"log:\n  level: debug\nwatch: \n  namespace: \n  interval: 30s"}}'
```

