#!/bin/bash
set -e -x -o pipefail

cd ../video_service && make upload
cd ../scheduler && make upload
cd ../user_service && make upload
cd ../front_api && make upload
cd ../webapp && make upload