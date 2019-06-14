SHELL=/bin/bash -o pipefail

GO ?= go
DOCKER ?= docker
HELM ?= helm

REGISTRY_IMAGE="codenotary/kube-notary:dev"

export GO111MODULE=on


.PHONY: kube-notary
kube-notary:
	GOOS=linux GOARCH=amd64 $(GO) build ./cmd/kube-notary

.PHONY: image
image:
	$(DOCKER) build -t $(REGISTRY_IMAGE) -f ./Dockerfile .

.PHONY: image.push
image.push: image
	$(DOCKER) push $(REGISTRY_IMAGE)

.PHONY: kubernetes/kube-notary
kubernetes/kube-notary:
	$(HELM) template -n kube-notary helm/kube-notary --output-dir ./kubernetes
	for f in ./kubernetes/kube-notary/templates/*; do grep -E "helm|Tiller" -v $$f > $$f.tmp; rm $$f; mv $$f.tmp $$f; done

.PHONY: test/e2e
test/e2e:
	cd ./test/e2e && ./run.sh
