SERVER_CONTAINER = infiniti-server
LINE_THROUGH = "=============================================================="

PKG := -d -v ./...
CMD := echo "Specify a command to run"

docker-build:
	@docker build -t infiniti-bank . && \
	if [ "${shell docker ps -aq -f name=^$(SERVER_CONTAINER)$}" ]; then \
		echo "$(LINE_THROUGH)\nContainer already exists, removing..."; \
		docker rm $(SERVER_CONTAINER); \
	fi; \
	echo "$(LINE_THROUGH)\nCreating container..."; \
	docker container create --name=$(SERVER_CONTAINER) --volume=$(PWD):/infiniti/ infiniti-bank && \
	docker start $(SERVER_CONTAINER) && \
	echo "$(LINE_THROUGH)\nInstalling dependencies..."; \
	make install && \
	echo "Done!\n$(LINE_THROUGH)"; \
	make stop

docker-logs:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	docker logs $(SERVER_CONTAINER)

start:
	@docker start $(SERVER_CONTAINER); \
		docker exec $(SERVER_CONTAINER) go build -o /infiniti/bin/infiniti /infiniti/cmd/main.go && \
	    echo "$(LINE_THROUGH)"; \
		docker exec $(SERVER_CONTAINER) /infiniti/bin/infiniti; \
		echo "$(LINE_THROUGH)"; \
		make stop

shell:
	@if [ -z "${shell docker ps -q -f name=^$(SERVER_CONTAINER)$}" ]; then \
		docker start $(SERVER_CONTAINER); \
	fi && \
	docker exec $(SERVER_CONTAINER) $(CMD)

stop:
	@echo "Stopping $(SERVER_CONTAINER)" \
		&& docker stop $(SERVER_CONTAINER)

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