#!/bin/bash
sed 's/localhost/nginx/g' -i ./webapp/.env
sed 's/localhost/nginx/g' -i ./nginx/nginx.conf
sed 's/localhost/nginx/g' -i ./.env.dev
sed 's/localhost/nginx/g' -i ./.env
