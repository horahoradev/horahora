# syntax=docker/dockerfile:1.2

# Using parts of https://github.com/six8/pyinstaller-alpine/blob/develop/python3.7.Dockerfile as reference
# For alpine, we need to build pyinstaller manually so its supported on alpine
FROM python:3.9-alpine as pyinstaller-builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

RUN apk add --update --no-cache zlib-dev musl-dev libc-dev libffi-dev gcc g++ git pwgen
RUN pip install pycrypto

RUN git clone --depth 1 --branch v4.5.1 https://github.com/pyinstaller/pyinstaller.git /pyinstaller
WORKDIR /pyinstaller
# install pyinstaller into a virtualenv (so we can copy the env to the other stages)
RUN python -m venv env
RUN . ./env/bin/activate && \
    cd bootloader && \
    CFLAGS="-Wno-stringop-overflow -Wno-stringop-truncation" python ./waf configure --no-lsb all && \
    pip install -r ../requirements.txt

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM python:3.9-alpine as ytdl-builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /

# download youtube-dl and prepare it for usage
RUN apk add --update --no-cache git gcc musl-dev
RUN git clone --depth 1 https://github.com/animelover1984/youtube-dl.git youtube-dl
WORKDIR /youtube-dl
RUN pip install -r requirements.txt
COPY --from=pyinstaller-builder /pyinstaller /pyinstaller
# we provide --paths here because for some reason pyinstaller doesnt figure that out on its own
RUN /pyinstaller/env/bin/python /pyinstaller/pyinstaller.py --paths /usr/local/lib/python3.9/site-packages --onefile -n youtube-dl youtube_dl/__main__.py

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM python:3.9-alpine as ytdlp-builder
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /

# download yt-dlp and prepare it for usage
RUN apk add --update --no-cache git gcc musl-dev
RUN git clone --depth 1 https://github.com/horahoradev/yt-dlp.git yt-dlp
WORKDIR /yt-dlp
RUN pip install -r requirements.txt
COPY --from=pyinstaller-builder /pyinstaller /pyinstaller
# we provide --paths here because for some reason pyinstaller doesnt figure that out on its own
RUN /pyinstaller/env/bin/python /pyinstaller/pyinstaller.py --paths /usr/local/lib/python3.9/site-packages --onefile -n yt-dlp yt_dlp/__main__.py

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

FROM alpine:3.14 AS runtime
LABEL org.opencontainers.image.source=https://github.com/horahoradev/horahora

WORKDIR /horahora/scheduler

COPY --from=ytdl-builder /youtube-dl/dist/youtube-dl /usr/local/bin/youtube-dl
COPY --from=ytdlp-builder /yt-dlp/dist/yt-dlp /usr/local/bin/yt-dlp
COPY --from=builder /scheduler.bin /scheduler.bin

ENTRYPOINT ["/scheduler.bin"]
