SHELL=/bin/bash -o pipefail

GO ?= go
DOCKER ?= docker

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
