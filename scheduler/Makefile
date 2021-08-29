all: proto docker

proto : ./protocol/scheduler.proto
	protoc -I=protocol scheduler.proto --go_out=plugins=grpc:protocol

docker : Dockerfile
	docker build -t scheduler .

upload : Dockerfile
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/scheduler .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/scheduler

build : Dockerfile
	eval $(minikube docker-env)
	docker build -t scheduler:latest .
