#!/bin/bash
set -euo pipefail

docker build -t grpcutil . --no-cache
docker run -v $(pwd)/video_service:/video_service -v $(pwd)/scheduler:/scheduler -v $(pwd)/user_service:/user_service -v $(pwd)/partyservice:/partyservice -it -t grpcutil /bin/bash
