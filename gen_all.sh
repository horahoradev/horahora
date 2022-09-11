#!/bin/bash
set -euo pipefail

$(cd video_service/protocol && protoc --go_out=. -I. -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=module=github.com/horahoradev/horahora/video_service/protocol  --validate_opt=module=github.com/horahoradev/horahora/video_service/protocol --validate_out="lang=go:."   videoservice.proto)
$(cd scheduler/protocol && protoc --go_out=. -I.  -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate --go_opt=paths=source_relative   --go-grpc_out=. --go-grpc_opt=module=github.com/horahoradev/horahora/scheduler/protocol  --validate_opt=module=github.com/horahoradev/horahora/scheduler/protocol --validate_out="lang=go:."   scheduler.proto)
$(cd user_service/protocol && protoc --go_out=. -I.  -I ${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=module=github.com/horahoradev/horahora/user_service/protocol  --validate_opt=module=github.com/horahoradev/horahora/user_service/protocol --validate_out="lang=go:."   userservice.proto)
