up:
	./generate.sh && \
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build --parallel --progress=plain && \
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f --tail=100
