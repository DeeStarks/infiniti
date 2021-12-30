# INFINITI

Inifiniti is a bank app that allows for the management of bank accounts and transactions. NB: This is an open source project meant for practice purposes.

- Architecture: Hexagonal
- Language: Go
- Containerization: Docker
- Database: Postgres
- Command Management: Makefile


### Setup

NB: Make sure to have docker installed and running.

- `git clone https://gihub.com/deestarks/infiniti.git`
- `cd infiniti`
- `make docker-build` # This will build the docker image
- `make start` # This will start the docker container.
- `make stop` # To stop the docker container.
- `make build` # Builds the application into the `bin/infiniti` binary.