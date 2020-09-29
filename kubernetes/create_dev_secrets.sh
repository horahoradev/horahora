#!/bin/bash
set -e -x -o pipefail

kubectl create secret generic videodb --from-literal=password=$1
kubectl create secret generic userdb --from-literal=password=$2
kubectl create secret generic scheduledb --from-literal=password=$3