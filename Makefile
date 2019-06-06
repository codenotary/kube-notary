SHELL=/bin/bash -o pipefail

GO ?= go
DOCKER ?= docker

REGISTRY_IMAGE="quay.io/leogr/kubewatch:dev"

export GO111MODULE=on


.PHONY: kubewatch
kubewatch:
	GOOS=linux GOARCH=amd64 $(GO) build ./cmd/kubewatch

.PHONY: image
image:
	$(DOCKER) build -t $(REGISTRY_IMAGE) -f ./Dockerfile .

.PHONY: image.push
image.push: image
	$(DOCKER) push $(REGISTRY_IMAGE)
