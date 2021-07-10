up:
	./generate.sh && \
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose up --build -d

down:
	docker-compose down

logs:
	docker-compose logs -f --tail=100
