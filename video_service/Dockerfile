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

FROM alpine:3.14 as mp4box-builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

# MP4Box is part of gpac
RUN apk add --update --no-cache git gcc make musl-dev zlib-dev zlib-static
RUN git clone --depth 1 https://github.com/gpac/gpac
WORKDIR /gpac
RUN ./configure --static-bin
RUN make -j
RUN cp bin/gcc/MP4Box /usr/local/bin/MP4Box

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM golang:1.16-alpine as builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/videoservice

RUN apk add --update --no-cache gcc musl-dev

# download modules
COPY go.mod /horahora/videoservice/
COPY go.sum /horahora/videoservice/

RUN go mod download

# build binary
COPY . /horahora/videoservice

RUN --mount=type=cache,target=/root/.cache/go-build go build -o /videoservice.bin

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM alpine:3.14 AS runtime
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/videoservice

RUN apk add --update --no-cache bash

COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=mp4box-builder /usr/local/bin/MP4Box /usr/local/bin/MP4Box
COPY --from=builder /videoservice.bin /videoservice.bin
COPY scripts/ /horahora/videoservice/scripts/

ENTRYPOINT ["/videoservice.bin"]