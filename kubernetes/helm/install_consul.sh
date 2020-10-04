#!/bin/bash
#helm repo add hashicorp https://helm.releases.hashicorp.com

helm install -f consul.yaml consul hashicorp/consul

# helm upgrade -f consul.yaml consul hashicorp/consul