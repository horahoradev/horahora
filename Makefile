up:
	./generate.sh && \
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build --parallel && \
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f --tail=100

proto:
	docker build -t grpcutil . && \
		docker run -it -t grpcutil
