# syntax=docker/dockerfile:1.2

# NOTE: because migrations come from outside the `front_api` directory,
#       we build this image from project root, symlinking this file to
#       Dockerfile.front_api

FROM golang:1.18.1-alpine as builder

WORKDIR /horahora/stomp_proxy

# build binary
COPY . /horahora/stomp_proxy

RUN go mod vendor && \
 go build -mod=vendor -o /stomp_proxy.bin

ENTRYPOINT ["/stomp_proxy.bin"]
