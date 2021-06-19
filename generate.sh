#!/bin/bash
env ./secrets.env

envsubst < docker-compose.yml.envs > docker-compose.yml