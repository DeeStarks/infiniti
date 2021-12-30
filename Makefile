SERVER_CONTAINER := infiniti-server
LINE_THROUGH := "=============================================================="

PKG := -d -v ./...

docker-build:
	@docker build -t infiniti-bank . && \
	if [ ! "$(docker ps -q -f name=^$(SERVER_CONTAINER)$)" ]; then \
		echo "$(LINE_THROUGH)\nContainer already exists, removing...\n$(LINE_THROUGH)"; \
		docker stop $(SERVER_CONTAINER) && \
		docker rm $(SERVER_CONTAINER); \
	fi; \
	echo "$(LINE_THROUGH)\nCreating container...\n$(LINE_THROUGH)"; \
	docker container create --name=$(SERVER_CONTAINER) --volume=$(PWD):/infiniti/ infiniti-bank && \
	docker start $(SERVER_CONTAINER) && \
	echo "$(LINE_THROUGH)\nInstalling dependencies\n$(LINE_THROUGH)"; \
	make install; \
	make stop

start:
	@docker start $(SERVER_CONTAINER); \
		docker exec $(SERVER_CONTAINER) go build -o /infiniti/bin/infiniti /infiniti/cmd/main.go && \
	    echo "$(LINE_THROUGH)"; \
		docker exec $(SERVER_CONTAINER) /infiniti/bin/infiniti; \
		echo "$(LINE_THROUGH)"; \
		make stop

stop:
	@echo "Stopping $(SERVER_CONTAINER)" \
		&& docker stop $(SERVER_CONTAINER)

install:
	@docker exec $(SERVER_CONTAINER) go get $(PKG)

build:
	@docker exec $(SERVER_CONTAINER) go build -o /infiniti/bin/infiniti /infiniti/cmd/main.go && \
		docker restart $(SERVER_CONTAINER)

tidy:
	@docker exec $(SERVER_CONTAINER) go mod tidy