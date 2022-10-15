#!/bin/bash
sed 's/nginx/localhost/g' -i ./webapp/.env
sed 's/nginx/localhost/g' -i ./nginx/nginx.conf
sed 's/nginx/localhost/g' -i ./.env.dev
sed 's/nginx/localhost/g' -i ./.env
