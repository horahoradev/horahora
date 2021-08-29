build : Dockerfile
	docker build -t frontend:latest .

upload : Dockerfile
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/frontend .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/frontend