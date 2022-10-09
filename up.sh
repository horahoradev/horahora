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

# copy example env file if it doesn't exist
if [ ! -f .env ]
then
  cp configs/.env.example .env
  # oh no no no no no
  sed -i '$d' .env
  echo -n 'JWT_KEYPAIR="' >> .env
  openssl genrsa 2048 >> .env
  echo '"' >> .env
fi

# docker compose by default reads `.env` file
# so no need to pass it as an option
docker-compose build --parallel --progress=plain && \
docker-compose up --build -d
