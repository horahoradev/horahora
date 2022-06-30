#!/bin/bash
set -euo pipefail

# copy example env file if it doesn't exist
if [ ! -f .env ]
then
  cp configs/.env.example .env
fi

# docker compose by default reads `.env` file
# so no need to pass it as an option
docker-compose build --progress=plain && \
docker-compose up -d
