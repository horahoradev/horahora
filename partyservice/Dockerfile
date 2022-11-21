# syntax=docker/dockerfile:1.2

FROM golang:1.18.1-alpine as builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/partyservice

RUN apk add --update --no-cache gcc musl-dev

# download modules
COPY go.mod /horahora/partyservice/
COPY go.sum /horahora/partyservice/

RUN go mod download

# build binary
COPY . /horahora/partyservice

RUN go build -o /partyservice.bin

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM alpine:3.14 AS runtime
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/partyservice

COPY --from=builder /partyservice.bin /partyservice.bin

ENTRYPOINT ["/partyservice.bin"]
