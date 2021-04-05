#!/bin/bash
set -e -x -o pipefail

kubectl create secret generic aws-key-id --from-file=$HOME/.aws/aws_access_key_id
kubectl create secret generic aws-secret-key --from-file=$HOME/.aws/aws_secret_access_key
kubectl create secret generic aws-region --from-file=$HOME/.aws/aws_region
