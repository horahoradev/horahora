#!/bin/bash

./generate.sh && \
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build --parallel --progress=plain && \
docker-compose up -d