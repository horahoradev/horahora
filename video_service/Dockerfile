FROM golang:1.13.7-buster

COPY . /videoservice

WORKDIR /videoservice

RUN apt install -y git
RUN go get github.com/go-delve/delve/cmd/dlv

RUN go build --mod=vendor -gcflags "all=-N -l" -o videoservice .

RUN apt update && \
    apt install -y ffmpeg

ENTRYPOINT ["./videoservice"]
