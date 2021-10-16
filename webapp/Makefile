all: proto docker

upload : Dockerfile
	docker build -t 908221837281.dkr.ecr.us-west-1.amazonaws.com/react_frontend .
	docker push 908221837281.dkr.ecr.us-west-1.amazonaws.com/react_frontend

build : Dockerfile
	eval $(minikube docker-env)
	docker build -t scheduler:latest .
