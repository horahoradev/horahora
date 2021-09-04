#!/bin/bash
set -euo pipefail

protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative    ./video_service/protocol/videoservice.proto
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative    ./user_service/protocol/userservice.proto
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative    ./scheduler/protocol/scheduler.proto