FROM golang:1.13.7-buster

COPY . /videoservice

WORKDIR /videoservice

RUN go build --mod=vendor -o videoservice .

RUN apt update && \
    apt install -y ffmpeg

ENTRYPOINT ["./videoservice"]
