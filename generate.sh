#!/bin/bash
set -x 


cd envsubst && sudo docker build -t envsubst . && cd ..

sudo docker run -v $(pwd)/:/envsubst envsubst /bin/bash -c "source ./envsubst/secrets.env.template && envsubst < envsubst/docker-compose.yml.envs > envsubst/docker-compose.yml"