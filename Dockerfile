# Copyright (c) 2019 vChain, Inc. All Rights Reserved.
# This software is released under GPL3.
# The full license information can be found under:
# https://www.gnu.org/licenses/gpl-3.0.en.html

FROM golang:1.15 as builder
WORKDIR /src

RUN apt-get install --no-install-recommends -y openssh-client
# Allow downloading vcn-enterprise using ssh agent forwarding
RUN mkdir ~/.ssh
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
ENV GOPRIVATE=github.com/codenotary/vcn-enterprise
RUN git config --global url."git@github.com:codenotary/vcn-enterprise".insteadOf "https://github.com/codenotary/vcn-enterprise"

COPY . .
RUN --mount=type=ssh \
  make kube-notary

FROM alpine:3.15

RUN apk update && apk upgrade && apk add ca-certificates curl musl && rm -rf /var/cache/apk/*

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN echo "curl -s 127.0.0.1:9581/results?output=bulk_sign" > /bin/bulk_sign \
    && chmod +x /bin/bulk_sign

COPY --from=builder /src/kube-notary /bin/kube-notary

RUN mkdir .vcn

ENTRYPOINT [ "/bin/kube-notary" ]
