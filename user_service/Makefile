all: proto

proto : ./protocol/userservice.proto
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative    protocol/userservice.proto 


build : Dockerfile
	eval $(minikube docker-env)
	docker build -t userservice:latest .

upload : Dockerfile
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/userservice .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/userservice


mock : ./protocol/userservice.pb.go
	cd protocol && mockgen -destination=mocks/mock_client.go -package=mocks . UserServiceClient