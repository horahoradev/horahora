#!/bin/bash
set -euo pipefail

# Run database setup after postgres starts (meaning, this runs in the background)
{
  until pg_isready -h localhost -U $POSTGRES_USER; do
    echo "Waiting for postgres to start..."
    sleep 1
  done

  # dodge the post-init restart lol
  sleep 20

  echo "Creating databases"
  # TODO(ivan): how do we want to make this so it runs only once
  psql --user=$POSTGRES_USER -a -f /postgres/create_db.sql

  echo "Running migrations"
  flyway -user=$POSTGRES_USER -password=$POSTGRES_PASSWORD -url=jdbc:postgresql://localhost:5432/userservice -locations=filesystem:///postgres/user_service/migrations migrate
  flyway -user=$POSTGRES_USER -password=$POSTGRES_PASSWORD -url=jdbc:postgresql://localhost:5432/scheduler -locations=filesystem:///postgres/scheduler/migrations migrate
  flyway -user=$POSTGRES_USER -password=$POSTGRES_PASSWORD -url=jdbc:postgresql://localhost:5432/videoservice -locations=filesystem:///postgres/video_service/migrations migrate

  echo "Seeding databases"
  # TODO(ivan): how do we want to make this so it runs only once
  psql --user=$POSTGRES_USER -a -f /postgres/seed_db.sql
} &

exec docker-entrypoint.sh postgres
