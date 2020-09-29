all: proto

proto : ./protocol/userservice.proto
	protoc -I=protocol userservice.proto --go_out=plugins=grpc:protocol


build : Dockerfile
	eval $(minikube docker-env)
	docker build -t userservice:latest .

upload : Dockerfile
	docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/userservice .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/userservice


mock : ./protocol/userservice.pb.go
	cd protocol && mockgen -destination=mocks/mock_client.go -package=mocks . UserServiceClient