#!/bin/bash
set -euo pipefail

if [ -f secrets.env.template ]
then
  echo "It looks like you are using old environment variables setup"
  echo "Read \"./docs/env-migrate.md\" on how to move into the new system."
  exit 1
fi

if [ -f docker-compose.yml.envs ]
then
  echo "It looks like you are using old environment variables setup"
  echo "Read \"./docs/env-migrate.md\" on how to move into the new system."
  exit 1
fi

sudo docker build -f Dockerfile.template -t gen . && docker run -it -v $(pwd):/gen -t gen

# docker compose by default reads `.env` file
# so no need to pass it as an option
DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 sudo docker compose build --parallel --progress=plain && \
sudo docker compose up -d
