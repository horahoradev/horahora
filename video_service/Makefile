all: proto docker

proto : ./protocol/videoservice.proto
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative    protocol/videoservice.proto 

docker : Dockerfile
	docker build -t videoservice .

upload : Dockerfile
	docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/videoservice .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/videoservice

mock : ./protocol/videoservice.pb.go
	cd protocol && mockgen -destination=mocks/mock_server.go -package=mocks . VideoService_UploadVideoServer

build : Dockerfile
	docker build -t videoservice:latest .
