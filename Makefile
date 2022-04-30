include .env # Load environment variables from .env
export

LINE_THROUGH = "=============================================================="

db-shell:
	@echo "$(LINE_THROUGH)\nLaunching database shell..."; \
		docker-compose exec infiniti-db psql -U $(DB_USER) -d $(DB_NAME); \
		echo "Database shell closed!\n$(LINE_THROUGH)"

start:
	@docker-compose up -d

logs:
	@docker-compose logs

keygen:
	@docker-compose exec infiniti-web openssl rand -hex 32 | xargs

TESTDIR := .
test:
	@echo "$(LINE_THROUGH)\nRunning tests..."; \
		docker-compose exec infiniti-web go test -v $(TESTDIR)/... && \
		echo "Tests complete!\n$(LINE_THROUGH)"

shell:
	@docker-compose exec infiniti-web /bin/bash

stop:
	@echo "Stopping Infiniti..."; \
		docker-compose down; \
		echo "Done!\n$(LINE_THROUGH)"

PKG := -d -v ./...
install:
	@docker-compose exec infiniti-web go get $(PKG)

build:
	@docker-compose exec infiniti-web go build -o /infiniti/bin/infiniti /infiniti/cmd/main.go && \
		docker-compose restart infiniti-web

tidy:
	@docker-compose exec infiniti-web go mod tidy

MIGRATION_NAME := $(shell date +%Y%m%d%H%M%S)
create-migrations:
	@docker-compose exec infiniti-web migrate create -ext sql -dir /infiniti/adapters/framework/db/migrations -seq $(MIGRATION_NAME)

migrate-up:
	@docker-compose exec infiniti-web migrate -path /infiniti/adapters/framework/db/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -verbose up

MIGRATE_VERSION := ""
migrate-down:
	@if [ -z "$(MIGRATE_VERSION)" ]; then \
		@docker-compose exec infiniti-web migrate -path /infiniti/adapters/framework/db/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -verbose down; \
	else \
		docker-compose exec infiniti-web migrate -path /infiniti/adapters/framework/db/migrations -database "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -verbose force $(MIGRATE_VERSION); \
	fi


# CLI app commands
cli-createadmin:
	@docker-compose exec infiniti-web ./bin/infiniti -c cli createadmin

cli-help:
	@docker-compose exec infiniti-web ./bin/infiniti -c cli help

cli-version:
	@docker-compose exec infiniti-web ./bin/infiniti -c cli version