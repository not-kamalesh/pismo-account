APP_NAME := pismo-account
DOCKER_IMAGE := pismo-account:latest

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           - Build Go binary locally"
	@echo "  run             - Run Go app locally"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container (needs MySQL separately)"
	@echo "  up              - docker-compose up (app + MySQL)"
	@echo "  down            - docker-compose down"
	@echo "  logs            - Show logs from docker-compose services"
	@echo "  ps              - List docker-compose services"

.PHONY: build
build:
	go build -o bin/$(APP_NAME) ./cmd

.PHONY: run
run:
	go run ./cmd

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run:
	docker run --rm \
		-p 8080:8080 \
		$(DOCKER_IMAGE)

.PHONY: up
up:
	docker compose up -d --build

.PHONY: down
down:
	docker compose down

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: ps
ps:
	docker compose ps
