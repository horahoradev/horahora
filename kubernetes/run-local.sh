#!/bin/bash
set -e -x -o pipefail
# This script will run the service locally

# 1. create services/deployments
kubectl apply -f local.yaml

# 2. Apply migrations
kill -9 $(ps -aux | grep "svc/scheduledb 5432:5432" | awk '{print $2}') || true
kubectl port-forward svc/scheduledb 5432:5432 &
sleep 20 # lol
psql --host=localhost -c 'create database scheduler' --user=guest || true
psql --host=localhost -c 'create database userservice' --user=guest || true
psql --host=localhost -c 'create database videoservice' --user=guest || true

proxyPID=$!
flyway -user=guest -password=guest -url=jdbc:postgresql://localhost:5432/scheduler -locations=filesystem://$(pwd)/../scheduler/migrations migrate
flyway -user=guest -password=guest -url=jdbc:postgresql://localhost:5432/userservice -locations=filesystem://$(pwd)/../user_service/migrations migrate
flyway -user=guest -password=guest -url=jdbc:postgresql://localhost:5432/videoservice -locations=filesystem://$(pwd)/../video_service/migrations migrate

sleep 30
./run-tests.sh

echo "Press enter to kill"
read
kill -9 $proxyPID