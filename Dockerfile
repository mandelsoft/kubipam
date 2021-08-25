#############      builder       #############
FROM golang:1.15.4 AS builder

ARG TARGETS=dev

WORKDIR /go/src/github.com/mandelsoft/kubelink
COPY . .

RUN make $TARGETS

############# base
FROM alpine:3.11.3 AS base

#############      command     #############
FROM base AS command

RUN apk add iptables
COPY --from=builder /go/bin/kubipam /kubipam

WORKDIR /

ENTRYPOINT ["/kubipam"]
