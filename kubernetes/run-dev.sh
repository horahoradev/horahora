#!/bin/bash
set -e -x -o pipefail
# This script will run the service on dev

# 1. create services/deployments
kubectl apply -f develop.yaml