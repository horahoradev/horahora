#!/bin/bash
set -e -x -o pipefail

flyway -user=userservice -password=$1 -url=jdbc:postgresql://userdb-dev.cwioxjkfbfkg.us-west-1.rds.amazonaws.com:5432/userservice -locations=filesystem://$(pwd)/../user_service/migrations migrate
flyway -user=scheduler -password=$2 -url=jdbc:postgresql://scheduledb-dev.cwioxjkfbfkg.us-west-1.rds.amazonaws.com:5432/scheduler -locations=filesystem://$(pwd)/../scheduler/migrations migrate
flyway -user=videoservice -password=$3 -url=jdbc:postgresql://videodb-dev.cwioxjkfbfkg.us-west-1.rds.amazonaws.com:5432/videoservice -locations=filesystem://$(pwd)/../video_service/migrations migrate


 https://horahora-dev-otomads.s3-us-west-1.amazonaws.com/542e0989-2952-11eb-aae3-9a2d98d52987258050498.mpd

