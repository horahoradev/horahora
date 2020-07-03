#!/bin/bash
set -e -x -o pipefail

# 1. Integration tests
kubectl exec $(kubectl get pods | grep "scheduler" | awk '{print $1}') -- /bin/bash -c "go test --mod=vendor -timeout=10m -tags=integration ./..."
kubectl exec $(kubectl get pods | grep "videoservice" | awk '{print $1}') -- /bin/bash -c "cd internal/models/ && go test ."
kubectl exec $(kubectl get pods | grep "userservice" | awk '{print $1}') -- /bin/bash -c "cd internal/auth/ && go test ."