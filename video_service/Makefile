all: proto docker

proto : ./protocol/videoservice.proto
	protoc -I=protocol videoservice.proto --go_out=plugins=grpc:protocol

docker : Dockerfile
	docker build -t videoservice .

mock : ./protocol/videoservice.pb.go
	cd protocol && mockgen -destination=mocks/mock_server.go -package=mocks . VideoService_UploadVideoServer

build : Dockerfile
	eval $(minikube docker-env)
	docker build -t videoservice:latest .
