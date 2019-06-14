# Manual installation

This folder contains templates to manually install `kube-notary` by using `kubectl`.

Edit files within [kube-notary/templates](kube-notary/templates) accordingly to your needs, then:

```
 kubectl apply \
    -f kube-notary/templates/serviceaccount.yaml \
    -f kube-notary/templates/clusterrolebinding.yaml \
    -f kube-notary/templates/configmap.yaml \
    -f kube-notary/templates/deployment.yaml \
    -f kube-notary/templates/service.yaml
```