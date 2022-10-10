# syntax=docker/dockerfile:1.2

FROM golang:1.18.1-alpine as builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/userservice

RUN apk add --update --no-cache gcc musl-dev

# download modules
COPY go.mod /horahora/userservice/
COPY go.sum /horahora/userservice/

RUN go mod download

# build binary
COPY . /horahora/userservice

RUN go build -o /userservice.bin

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM alpine:3.14 AS runtime
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/userservice

COPY --from=builder /userservice.bin /userservice.bin

ENTRYPOINT ["/userservice.bin"]
