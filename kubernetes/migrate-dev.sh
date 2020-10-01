#!/bin/bash
set -e -x -o pipefail

flyway -user=userservice -password=$1 -url=jdbc:postgresql://userdb-dev.cwioxjkfbfkg.us-west-1.rds.amazonaws.com:5432/userservice -locations=filesystem://$(pwd)/../user_service/migrations migrate
flyway -user=scheduler -password=$2 -url=jdbc:postgresql://scheduledb-dev.cwioxjkfbfkg.us-west-1.rds.amazonaws.com:5432/scheduler -locations=filesystem://$(pwd)/../scheduler/migrations migrate
flyway -user=videoservice -password=$3 -url=jdbc:postgresql://videodb-dev.cwioxjkfbfkg.us-west-1.rds.amazonaws.com:5432/videoservice -locations=filesystem://$(pwd)/../video_service/migrations migrate
