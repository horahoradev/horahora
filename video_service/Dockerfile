FROM golang:1.14.6-alpine3.12

COPY . /videoservice

WORKDIR /videoservice

RUN apk add git
RUN apk add bash

RUN go get github.com/go-delve/delve/cmd/dlv

RUN cgo_enabled=0 go build --mod=vendor -o videoservice .

RUN apk add ffmpeg

ENTRYPOINT ["./videoservice"]
