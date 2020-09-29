all: proto docker

proto : ./protocol/videoservice.proto
	protoc -I=protocol videoservice.proto --go_out=plugins=grpc:protocol

docker : Dockerfile
	docker build -t videoservice .

upload : Dockerfile
	docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/videoservice .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/videoservice

mock : ./protocol/videoservice.pb.go
	cd protocol && mockgen -destination=mocks/mock_server.go -package=mocks . VideoService_UploadVideoServer

build : Dockerfile
	docker build -t videoservice:latest .
