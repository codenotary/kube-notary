SHELL=/bin/bash -o pipefail

GO ?= go
DOCKER ?= docker
HELM ?= helm

REGISTRY_IMAGE="codenotary/kube-notary:latest"
TEST_FLAGS ?= -v -race

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

.PHONY: kubernetes
kubernetes:
	rm -rf kubernetes/kube-notary
	rm -rf kubernetes/kube-notary-namespaced
	$(HELM) template -n kube-notary helm/kube-notary --set watch.namespace="default" --output-dir ./kubernetes
	for f in ./kubernetes/kube-notary/templates/*; do grep -E "helm|Tiller" -v $$f > $$f.tmp; rm $$f; mv $$f.tmp $$f; done
	mv kubernetes/kube-notary kubernetes/kube-notary-namespaced
	$(HELM) template -n kube-notary helm/kube-notary --output-dir ./kubernetes
	for f in ./kubernetes/kube-notary/templates/*; do grep -E "helm|Tiller" -v $$f > $$f.tmp; rm $$f; mv $$f.tmp $$f; done

.PHONY: CHANGELOG.md
CHANGELOG.md:
	git-chglog -o CHANGELOG.md

.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test ${TEST_FLAGS} ./...

.PHONY: test/e2e.local
test/e2e.local:
	$(DOCKER) build -t kube-notary:test -f ./Dockerfile .
	cd ./test/e2e && ./run.sh

.PHONY: test/e2e
test/e2e:
	$(DOCKER) build -t kube-notary:test -f ./Dockerfile .
	$(DOCKER) build -t kube-notary-test-e2e -f Dockerfile.test-e2e .
	$(DOCKER) run --rm -v "/var/run/docker.sock:/var/run/docker.sock:ro" --network host kube-notary-test-e2e

