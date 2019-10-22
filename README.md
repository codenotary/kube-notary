# kube-notary
> A Kubernetes watchdog for verifying image trust with CodeNotary.

![alt text](https://raw.githubusercontent.com/vchain-us/kube-notary/master/docs/images/info.png "CodeNotary for Kubernetes")

## How it works

**kube-notary** is a monitoring tool for *Continuous Verification* (CV) via [CodeNotary](https://codenotary.io). 
The idea behind CV is to continuously monitor your cluster at runtime and be notified when unknown or untrusted container images are running.

Once `kube-notary` is installed within your cluster, all pods are checked every minute (interval and other settings can be [configured](#Configuration)). 
For each of the running containers in each pod, `kube-notary` resolves the `ImageID` of the container's image to the actual image's hash and finally looks up the [hash's signature in the CodeNotary's blockchain](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/signatures.md#signatures).

Furthermore, kube-notary provides a built-in exporter for sending verification [metrics](#Metrics) to Prometheus, which can then that can be easily visualized with the provided [grafana dashboard](grafana). 

Images you trust can be signed by using the CodeNotary [vcn](https://github.com/vchain-us/vcn) CLI tool.

https://infograph.venngage.com/ps/ex4ECrROPCQ/codenotary-for-kubernetes


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
> See the [Configuration](#Configuration) paragraph for detailed instructions.

### Namespaced

If you do not have cluster-wide access, you can still install `kube-notary` within a single namespace, using:
```
helm install -n kube-notary helm/kube-notary --set watch.namespace="default"
```

When so configured, a namespaced `Role` will be created instead of the default `ClusterRole` to accommodate Kubernetes [RBAC](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) for a single namespace. `kube-notary` will get permission for, and will watch, the configured namespace only.

### Manual installation (without Helm)
Alternatively, it is possible to manually install `kube-notary` without using Helm. Instructions and templates for manual installation are within the [kubernetes folder](kubernetes).

## Uninstall

You can uninstall `kube-notary` at any time using:
```
helm delete --purge kube-notary
```

## Usage

`kube-notary` provides both detailed log output and a Prometheus metrics endpoint to monitor the verification status of your running containers. After the installation you will find instructions on how to get them.

Examples:
```
  # Metrics endpoint
  export SERVICE_NAME=service/$(kubectl get service --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
  echo "Check the metrics endpoint at http://127.0.0.1:9581/metrics"
  kubectl port-forward --namespace default $SERVICE_NAME 9581

  # Results endpoint
  export SERVICE_NAME=service/$(kubectl get service --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
  echo "Check the verification results endpoint at http://127.0.0.1:9581/results"
  kubectl port-forward --namespace default $SERVICE_NAME 9581

  # Stream logs
  export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
  kubectl logs --namespace default -f $POD_NAME

  # Bulk sign all running images
  export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
  kubectl exec --namespace default -t $POD_NAME sh /bin/bulk_sign > vcn_bulk_sign.sh
  chmod +x vcn_bulk_sign.sh && ./vcn_bulk_sign.sh
```
### Status page

`kube-notary` provides an embedded status page with details of running containers and verifications. Once `kube-notary` is up and running the page is available at `http://$SERVICE_ADDRESS:9581/status/` (the service address depends on your installation). If you installed it by using the provided helm chart, the easiest way to get the status page URL is:

```
  # Status page
  export SERVICE_NAME=service/$(kubectl get service --namespace default -l "app.kubernetes.io/name=kube-notary,app.kubernetes.io/instance=kube-notary" -o jsonpath="{.items[0].metadata.name}")
  echo "Status page at http://127.0.0.1:9581/status/"
  kubectl port-forward --namespace default $SERVICE_NAME 9581
``` 

### Metrics

If a Prometheus installation is running within your cluster, metrics provided by `kube-notary` will be automatically discovered. 
Furthermore, you can find an example of a preconfigured Grafana dashboard [here](grafana/).

## Configuration

By default, `kube-notary` is installed into the current namespace (you can change it by using `helm install --namespace`) but it will watch to pods in all namespaces.

At install time you can change any values of [helm/kube-notary/values.yaml](helm/kube-notary/values.yaml) by using the Helm's `--set` option.
For example, to instruct `kube-notary` to check only the `kube-system` namespace, just use:
```
helm install -n kube-notary helm/kube-notary --set watch.namespace="kube-system"
```

### Trusted signers

Choosing signers who can be trusted is easy. At install time, you can `--set` one of the following options.

#### Configure a list of trusted keys

Lists can be expressed by enclosing keys in `{` and `}` separated by `,`. For example:
```
helm install -n kube-notary helm/kube-notary --set "trust.keys={0x123..., 0x123...}"
```

#### Configure a trusted organization

```
helm install -n kube-notary helm/kube-notary --set "trust.org=your.organization.com"
```

#### Note 
If both `trust.org` and `trust.keys` are set, only `trust.org` will be used. 

If none is set, the last [signature](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/signatures.md#signatures) by the signer with the highest available [level](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/signatures.md#level) will be used during the verification.

### Runtime configuration

The following options within [helm/kube-notary/values.yaml](helm/kube-notary/values.yaml) have effect on the `kube-notary` runtime behavior.
```
# Runtime config
log:
  level: info # verbosity level, one of: trace, debug, info, warn, error, fatal, panic
watch: 
  namespace: "" # the namespace name to watch 
  interval: 60s # duration of the watching interval
trust:
  org: "" # ID of the trusted organization, if not empty `keys` above will be ignored
  keys: # array of signing keys to verify against
   - ...
   - ...
```

During the installation, they are stored in a `configmap`. Configuration hot-reloading is supported, so you can modify and apply the configmap while `kube-notary` is running. 

For example, to change the watching interval from default to `30s`:
```
kubectl patch configmaps/kube-notary \
    --type merge \
    -p '{"data":{"config.yaml":"log:\n  level: debug\nwatch: \n  namespace: \n  interval: 30s"}}'
```

## FAQ

### Why *Continuous Verification* ?

Things change over time. Suppose you signed an image because you trust it. Later, you find a security issue within the image or you just want to deprecate that version. When that happens you can simply use [vcn](https://github.com/vchain-us/vcn#basic-usage) to [untrust or unsupport](https://github.com/vchain-us/vcn/blob/master/docs/user-guide/signatures.md#statuses) that image version. Once the image is not trusted anymore, 
thanks to `kube-notary` you can easily discover if the image is still running somewhere in your cluster.

In general, verifying an image just before its execution is not enough because the image's status or the image that's used by a container can change over time. *Continuous Verification* ensures that you will always get noticed if an unwanted behavior occurs.

### How can I sign my image?

You can easily sign your container's images by using the [vcn CLI](https://github.com/vchain-us/vcn) we provide separately.

`vcn` supports local Docker installations out of the box using `docker://` prefix, followed by the image name or image reference. 
You just have to pull the image you want to sign, then finally run `vcn sign`. Detailed instructions can be found [here](https://github.com/vchain-us/vcn/blob/master/docs/DOCKERINTEGRATION.md).

Furthermore, if you want to bulk sign all images running inside your cluster, you will find below instructions to generate a script that automates the process.

Export `POD_NAME` setting it to the `kube-notary`'s pod name, then run:
```
kubectl exec --namespace default -t $POD_NAME sh /bin/bulk_sign > vcn_bulk_sign.sh
chmod +x vcn_bulk_sign.sh && ./vcn_bulk_sign.sh
```
> Note that a [CodeNotary](https://codenotary.io) account and a local installation of [vcn](https://github.com/vchain-us/vcn) are needed.
> Also, make sure your `kubectl` is pointing to the context you want to use. 

### How can I be notified when untrusted images are runnig?

First, Prometheus and Grafana need to be installed in your cluster.

Then it's easy to [create alerts](grafana#creating-alerts) using the provided [Grafana dashboard](grafana)

### Why my image cannot be signed? (manifest v2 schema 1)

The image manifest v2 schema 1 format is [deprecated](https://docs.docker.com/engine/deprecated/#pushing-and-pulling-with-image-manifest-v2-schema-1) in favor of the v2 schema 2 format.

Please, move to v2 schema 2 as soon as possible. Usually, you can fix that simply by pushing your image again on the registry.

### Cannot create resource "clusterrolebindings"

Recent versions of Kubernetes employ a [role-based access control](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) (or RBAC) system to drive authorization decisions. It might be possible that your account does not have enough privileges to create the `ClusterRole` needed to get cluster-wide access.

Please use a high privileged account to install `kube-notary`. Alternatively, if you don't have cluster-wide access, you can still install `kube-notary` to work in a single namespace which you can access. See the [namespaced installation](#Namespaced) paragraph for further details.

### Helm error: release kube-watch failed: namespaces "..." is forbidden
It might be possible that `tiller` (the Helm's server-side component) does not have permission to install `kube-notary`. 

When working within a [role-based access control](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) enabled Kubernetes installation, you may need to add a [service account with cluster-admin role](https://helm.sh/docs/using_helm/#tiller-and-role-based-access-control) for `tiller`.

The easier way to do that is just to create a `rbac-config.yaml` copying and pasting the [provided example in the Helm documentation](https://helm.sh/docs/using_helm/#example-service-account-with-cluster-admin-role), then:

```
$ kubectl create -f rbac-config.yaml
serviceaccount "tiller" created
clusterrolebinding "tiller" created
$ helm init --service-account tiller --history-max 200
```

## Testing
```
make test/e2e
```

## License

This software is released under [GPL3](https://www.gnu.org/licenses/gpl-3.0.en.html).
