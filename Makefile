docker-build:
	@docker build -t infiniti-bank .

start:
	@echo "=========================================================="; \
		docker run --name infiniti-server --volume=${PWD}:/infiniti/ infiniti-bank; \
		echo "=========================================================="; \
		make stop

stop:
	@echo "Stopping infiniti-server" \
		&& docker stop infiniti-server \
		&& docker rm infiniti-server

build:
	@docker exec -it infiniti-server go build -o /infiniti/bin/main /infiniti/cmd/main.go