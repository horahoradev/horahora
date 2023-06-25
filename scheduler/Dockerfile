# syntax=docker/dockerfile:1.2

FROM alpine:3.14 as ffmpeg-builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

RUN apk add --update --no-cache wget
# RUN wget https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz
# TODO(ivan): Using my server temporarily because johnvansickle's site is very slow
RUN wget https://media.sq10.net/ffmpeg-git-amd64-static.tar.xz
RUN tar -xvf ffmpeg-git-amd64-static.tar.xz
RUN cd ffmpeg-git-*-amd64-static && cp ffmpeg /usr/local/bin/ffmpeg

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM golang:1.18.1-alpine as builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/scheduler

RUN apk add --update --no-cache gcc musl-dev

# download modules
COPY go.mod /horahora/scheduler/
COPY go.sum /horahora/scheduler/

RUN go mod download

# build binary
COPY . /horahora/scheduler

RUN go build -o /scheduler.bin

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM python:3.9-alpine AS runtime
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

RUN apk add --update --no-cache zlib-dev musl-dev libc-dev libffi-dev gcc g++ git pwgen jq && git clone --depth 1 https://github.com/horahoradev/yt-dlp.git yt-dlp

# download yt-dlp and prepare it for usage
WORKDIR /yt-dlp
RUN pip install -r requirements.txt && ln -s /yt-dlp/yt-dlp.sh /usr/local/bin/yt-dlp

WORKDIR /horahora/scheduler

COPY --from=builder /scheduler.bin /scheduler.bin
COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg

ENTRYPOINT ["/scheduler.bin"]
