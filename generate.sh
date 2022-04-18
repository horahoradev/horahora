#!/bin/bash
set -x 

# If we don't want to use the default keypair, check to see if custom keypair has already been generated
# If so, don't generate a new one.
if [ "$USE_DEFAULT_KEYPAIR" = false ] ; then
    `test -f ~/.ssh/horahora.pem || openssl genrsa -out ~/.ssh/horahora.pem 2048`
fi

source ./secrets.env.template

envsubst < docker-compose.yml.envs > docker-compose.yml