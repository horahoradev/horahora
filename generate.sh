#!/bin/bash
source ./secrets.env

envsubst < docker-compose.yml.envs > docker-compose.yml