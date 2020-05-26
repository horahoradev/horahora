#!/bin/bash
set -e -x -o pipefail

eval $(minikube docker-env)


cd ../video_service && make build
cd ../user_service && make build
cd ../scheduler && make build