SERVER_CONTAINER = infiniti-server
DB_CONTAINER = infiniti-db
LINE_THROUGH = "=============================================================="

PKG := -d -v ./...

DB_NAME := infiniti
DB_USER := postgres
DB_PASS := infiniti

MIGRATION_NAME := $(shell date +%Y%m%d%H%M%S)

docker-build:
	@docker build -t infiniti-bank . ; \
	if [ "${shell docker ps -aq -f name=^$(SERVER_CONTAINER)$}" ]; then \
		echo "$(LINE_THROUGH)\nContainer already exists, removing..."; \
		docker rm $(SERVER_CONTAINER); \
	fi; \
	echo "$(LINE_THROUGH)\nCreating container..."; \
	docker container create --name=$(SERVER_CONTAINER) --volume=$(PWD):/infiniti/ infiniti-bank && \
	echo "$(LINE_THROUGH)\nInstalling dependencies..."; \
	make install; \
	make db-setup && \
	echo "Done!\n$(LINE_THROUGH)"; \
	make stop

docker-logs:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	docker logs $(SERVER_CONTAINER)

db-setup:
	@if [ -z "${shell docker ps -a -q -f name=^$(DB_CONTAINER)$}" ]; then \
		echo "$(LINE_THROUGH)\nSetting up database..."; \
		docker pull postgres:12.2-alpine && \
		docker run --name=$(DB_CONTAINER) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASS) -e POSTGRES_DB=$(DB_NAME) -d postgres:12.2-alpine && \
		docker stop infiniti-db && \
		echo "Database setup complete!\n$(LINE_THROUGH)"; \
	else \
		echo "$(LINE_THROUGH)\nDatabase already exists, skipping setup...\n$(LINE_THROUGH)"; \
	fi

db-shell:
	@echo "$(LINE_THROUGH)\nOpening database shell..."; \
		if [ -z "${shell docker ps -q -f name=^$(DB_CONTAINER)$}" ]; then \
			docker start $(DB_CONTAINER) && \
			docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME); \
		else \
			docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME); \
		fi && \
		docker stop $(DB_CONTAINER) && \
		echo "Database shell closed!\n$(LINE_THROUGH)"

start:
	@docker start $(DB_CONTAINER) && \
		docker start $(SERVER_CONTAINER) && \
		docker exec $(SERVER_CONTAINER) go build -o /infiniti/bin/infiniti /infiniti/cmd/main.go && \
	    echo "$(LINE_THROUGH)"; \
		make tidy && \
		docker exec $(SERVER_CONTAINER) /infiniti/bin/infiniti; \
		echo "$(LINE_THROUGH)"; \
		make stop

shell:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	docker exec -it $(SERVER_CONTAINER) /bin/bash

stop:
	@echo "Stopping $(SERVER_CONTAINER)..." \
		&& docker stop $(SERVER_CONTAINER); \
		echo "Stopping $(DB_CONTAINER)..." \
		&& docker stop $(DB_CONTAINER); \
		echo "Done!\n$(LINE_THROUGH)"

install:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	docker exec $(SERVER_CONTAINER) go get $(PKG)

build:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	docker exec $(SERVER_CONTAINER) go build -o /infiniti/bin/infiniti /infiniti/cmd/main.go && \
		docker restart $(SERVER_CONTAINER)

tidy:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	docker exec $(SERVER_CONTAINER) go mod tidy

create-migrations:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	docker exec $(SERVER_CONTAINER) migrate create -ext sql -dir /infiniti/adapters/framework/db/migrations -seq $(MIGRATION_NAME)

migrate-up:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	if [ -z "${shell docker ps -q -f name=^$(DB_CONTAINER)$}" ]; then \
		docker start $(DB_CONTAINER); \
	fi && \
	docker exec $(SERVER_CONTAINER) migrate -path /infiniti/adapters/framework/db/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@host.docker.internal:5432/$(DB_NAME)?sslmode=disable" -verbose up

migrate-down:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	if [ -z "${shell docker ps -q -f name=^$(DB_CONTAINER)$}" ]; then \
		docker start $(DB_CONTAINER); \
	fi && \
	docker exec -it $(SERVER_CONTAINER) migrate -path /infiniti/adapters/framework/db/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@host.docker.internal:5432/$(DB_NAME)?sslmode=disable" -verbose down