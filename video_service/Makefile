all: proto docker

proto : ./protocol/videoservice.proto
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative    protocol/videoservice.proto 

dockerproto : ./protocol/videoservice.proto
	docker build https://github.com/envoyproxy/protoc-gen-validate.git#main -t videoservicegrpc && \
	sudo docker run --entrypoint=/bin/bash --mount type=bind,source=${pwd}/protocol,target=/proto videoservicegrpc -c "cd /proto &&  make proto"

docker : Dockerfile
	docker build -t videoservice .

upload : Dockerfile
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/videoservice .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/videoservice

mock : ./protocol/videoservice.pb.go
	cd protocol && mockgen -destination=mocks/mock_server.go -package=mocks . VideoService_UploadVideoServer

build : Dockerfile
	docker build -t videoservice:latest .
