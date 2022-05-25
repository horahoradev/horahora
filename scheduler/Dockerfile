# syntax=docker/dockerfile:1.2

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM golang:1.16-alpine as builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/scheduler

RUN apk add --update --no-cache gcc musl-dev

# download modules
COPY go.mod /horahora/scheduler/
COPY go.sum /horahora/scheduler/

RUN go mod download

# build binary
COPY . /horahora/scheduler

RUN --mount=type=cache,target=/root/.cache/go-build go build -o /scheduler.bin

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM python:3.9-alpine AS runtime
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

RUN apk add --update --no-cache zlib-dev musl-dev libc-dev libffi-dev gcc g++ git pwgen && git clone --depth 2 https://github.com/horahoradev/yt-dlp.git yt-dlp

# download yt-dlp and prepare it for usage
WORKDIR /yt-dlp
RUN pip install -r requirements.txt && ln -s /yt-dlp/yt-dlp.sh /usr/local/bin/yt-dlp

WORKDIR /horahora/scheduler

COPY --from=builder /scheduler.bin /scheduler.bin

ENTRYPOINT ["/scheduler.bin"]
