all: proto docker

proto : ./protocol/scheduler.proto
	protoc -I=protocol scheduler.proto --go_out=plugins=grpc:protocol

docker : Dockerfile
	docker build -t scheduler .

build : Dockerfile
	eval $(minikube docker-env)
	docker build -t scheduler:latest .
