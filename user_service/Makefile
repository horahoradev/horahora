all: proto

proto : ./protocol/userservice.proto
	protoc -I=protocol userservice.proto --go_out=plugins=grpc:protocol


build : Dockerfile
	eval $(minikube docker-env)
	docker build -t userservice:latest .

mock : ./protocol/userservice.pb.go
	cd protocol && mockgen -destination=mocks/mock_client.go -package=mocks . UserServiceClient