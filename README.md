# INFINITI

Inifiniti is a bank app that allows for the management of bank accounts and transactions. NB: This is an open source project meant for practice purposes.

- Architecture: Hexagonal
- Language: Go
- Containerization: Docker
- Database: Postgres
- Command Management: Makefile


### Setup/Usage

NB: Make sure to have docker installed and running. For docker installation, see [here](https://docs.docker.com/get-docker/).

- `git clone https://gihub.com/deestarks/infiniti.git`
- `cd infiniti`
- `make docker-build` # This will build the docker image
- `make db-setup` # This will create the database with `DB_NAME=infiniti DB_USER=postgres DB_PASS=infiniti`. To change these values, pass them as arguments to the command (e.g. `make db-setup DB_NAME=example_db DB_USER=example_user DB_PASS=example_password`). NB: You might not need to run this command because the `make docker-build` command will create the database automatically. And if you need to specify the database name, user, and password, you can also do so by passing them as arguments to the command (e.g. `make docker-build DB_NAME=example_db DB_USER=example_user DB_PASS=example_password`).
- `make create-migrations` # This will create the migrations files for up and down in the `adapters/db/migrations`. Default name is a timestamp in format `YYYYMMDDHHMMSS`. To specify a different name, pass it as an argument to the command (e.g. `make create-migrations MIGRATION_NAME=example_name`).
- `make migrate-up` # This will run all `*.up.sql` migration files sequentially. NB: Pass the database credentials `DB_NAME, DB_USER, DB_PASS` as arguments if using a custom database name, user, or password.
- `make migrate-down` # Same as `make migrate-up` but runs `*.down.sql` files in reverse order.
- `make start` # This will start the docker container.
- `make stop` # To stop the docker container.
- `make build` # Builds the application into the `bin/infiniti` binary.
- `make install` # Installs all used dependencies in the application.
- `make install PKG=<dependency>` # Installs a specific dependency. (e.g. `make install PKG=github.com/gorilla/mux`)
- `make tidy` # Equivalent to `go mod tidy`, but in the docker container.
- `make shell` # Opens a shell in the docker container.