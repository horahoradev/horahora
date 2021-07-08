FROM golang:1.15.8-buster

COPY . /videoservice
WORKDIR /videoservice

#
#RUN go get github.com/go-delve/delve/cmd/dlv

RUN go mod vendor && \
    go build --mod=vendor -o videoservice .

RUN apt-get update && \
    apt-get install -y ffmpeg curl

COPY . /videoservice

RUN cgo_enabled=0 go build -o videoservice .

ENTRYPOINT ["./videoservice"]
