FROM golang:1.15.2-alpine3.12

RUN mkdir -p /videoservice/test_files
WORKDIR /videoservice

#RUN apk add git
RUN apk add bash
#
#RUN go get github.com/go-delve/delve/cmd/dlv

RUN apk add ffmpeg

COPY . /videoservice

RUN cgo_enabled=0 go build -o videoservice .

ENTRYPOINT ["./videoservice"]
