# syntax=docker/dockerfile:1.2

# NOTE: because migrations come from outside the `front_api` directory,
#       we build this image from project root, symlinking this file to
#       Dockerfile.front_api

FROM golang:1.18.1-alpine as builder

WORKDIR /horahora/front_api

# build binary
COPY . /horahora/front_api

RUN go mod vendor && \
 go build -mod=vendor -o /front_api.bin

ENTRYPOINT ["/front_api.bin"]
