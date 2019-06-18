# Manual installation

This folder contains templates to manually install `kube-notary` by using `kubectl`.
Before proceeding, make sure your local config is pointing to the context you want to use (ie. check `kubectl config current-context`).

## Cluster-wide

Edit files within [kube-notary/templates](kube-notary/templates) accordingly to your needs, then:

```
 kubectl apply \
    -f kube-notary/templates/serviceaccount.yaml \
    -f kube-notary/templates/role.yaml \
    -f kube-notary/templates/rolebinding.yaml \
    -f kube-notary/templates/configmap.yaml \
    -f kube-notary/templates/deployment.yaml \
    -f kube-notary/templates/service.yaml
```

## Namespaced

Files within [kube-notary-namespaced/templates](kube-notary-namespaced/templates) are compiled to work with `namespace: default`. If you are using a different namespace, please modify any instance of `namespace: default` before applying them.

Finally:
```
 kubectl apply \
    -f kube-notary-namespaced/templates/serviceaccount.yaml \
    -f kube-notary-namespaced/templates/role.yaml \
    -f kube-notary-namespaced/templates/rolebinding.yaml \
    -f kube-notary-namespaced/templates/configmap.yaml \
    -f kube-notary-namespaced/templates/deployment.yaml \
    -f kube-notary-namespaced/templates/service.yaml
```