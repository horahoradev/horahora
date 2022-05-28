all: proto docker

docker : Dockerfile
	docker build -t front_api .

upload : Dockerfile
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/frontend .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/frontend

build : Dockerfile
	eval $(minikube docker-env)
	docker build -t scheduler:latest .
