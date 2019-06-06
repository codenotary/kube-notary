# kubewatch
> A Kubernetes watchdog for verifying image trust with CodeNotary.

*Work in progress!*

## Usage

First, make sure `kubectl` is pointing to the context you want to use.

If you have RBAC enabled on your cluster, use the provided [kubernetes/clusterrolebinding.yaml](kubernetes/clusterrolebinding.yaml) to create role binding which will grant the default service account view permissions. Eventually modify it accoring to your needs. 

```
kubectl apply -f kubernetes/clusterrolebinding.yaml
```

Then create the configmap:
```
kubectl apply -f kubernetes/clusterrolebinding.yaml
```
> `kubewatch` supports configuration hot-reloading. You can modify and apply the configmap while `kubewatch` is running.

Finally, deploy `kubewatch`:
```
kubectl apply -f kubernetes/deployment.yaml
```

*TODO: provide Helm chart*

### Watch the log
```
kubectl logs -f deployment.apps/kubewatch-deployment
```

