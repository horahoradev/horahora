#!/bin/bash
source ./secrets.env.template

envsubst < docker-compose.yml.envs > docker-compose.yml