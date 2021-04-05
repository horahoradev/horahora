#!/bin/bash

mkdir ../data

# 1. DB creation
docker exec $(docker ps | grep "postgres" | awk '{print $1}') /bin/bash -c "psql --user=admin -a -f /sql/create_db.sql"

# 2. Schema migration
flyway -user=admin -password=password -url=jdbc:postgresql://localhost:5432/userservice -locations=filesystem://$(pwd)/user_service/migrations migrate
flyway -user=admin -password=password -url=jdbc:postgresql://localhost/scheduler -locations=filesystem://$(pwd)/scheduler/migrations migrate
flyway -user=admin -password=password -url=jdbc:postgresql://localhost:5432/videoservice -locations=filesystem://$(pwd)/video_service/migrations migrate

chmod 777 data