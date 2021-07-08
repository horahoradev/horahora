# syntax=docker/dockerfile:1.2

# NOTE: because migrations come from outside the `front_api` directory,
#       we build this image from project root, symlinking this file to
#       Dockerfile.front_api

FROM golang:1.16.5-buster as builder

WORKDIR /horahora/front_api

# download modules
COPY front_api/go.mod /horahora/front_api/
COPY front_api/go.sum /horahora/front_api/
COPY scheduler/go.mod /horahora/scheduler/
COPY scheduler/go.sum /horahora/scheduler/
COPY video_service/go.mod /horahora/video_service/
COPY video_service/go.sum /horahora/video_service/
COPY user_service/go.mod /horahora/user_service/
COPY user_service/go.sum /horahora/user_service/

RUN go mod download

# build binary
COPY front_api /horahora/front_api
COPY scheduler/protocol /horahora/scheduler/protocol
COPY video_service/protocol /horahora/video_service/protocol
COPY user_service/protocol /horahora/user_service/protocol

RUN --mount=type=cache,target=/root/.cache/go-build go build -o /front_api.bin

ENTRYPOINT ["/front_api.bin"]
