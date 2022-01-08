# INFINITI

Infiniti is a bank application that follows the hexagonal architecture, to design a system for managing accounts and user transactions. NB: This is for practice purposes only.

### Setup

NB: Make sure to have docker installed and running. For docker installation, see [here](https://docs.docker.com/get-docker/).

- `git clone https://gihub.com/deestarks/infiniti.git`
- `cd infiniti`
- `make docker-build` - To build the docker image
- `make migrate-up` - Migrates the database schema. For more info about this command, see no. 4 in the [List of Commands](#commands) section.
- `make start` - To start the application.


### Commands
1. `make docker-build` - Builds the docker image, and automatically installs all dependencies and sets up the database.
2. `make db-setup` - Creates the database with `DB_NAME=infiniti DB_USER=postgres DB_PASS=infiniti`. To change these values, pass them as arguments to the command (e.g. `make db-setup DB_NAME=example_db DB_USER=example_user DB_PASS=example_password`). NB: You might not need to run this command because the `make docker-build` command will create the database automatically. And if you need to specify the database name, user, and password, you can also do so by passing them as arguments to the command (e.g. `make docker-build DB_NAME=example_db DB_USER=example_user DB_PASS=example_password`).
3. `make create-migrations` - Creates the migrations files for up and down in the `adapters/db/migrations`. Default name is a timestamp in format `YYYYMMDDHHMMSS`. To specify a different name, pass it as an argument to the command (e.g. `make create-migrations MIGRATION_NAME=example_name`).
4. `make migrate-up` - Runs all `*.up.sql` migration files sequentially. NB: Pass the database credentials `DB_NAME, DB_USER, DB_PASS` as arguments if using a custom database name, user, or password.
5. `make migrate-down` - Same as `make migrate-up` but runs `*.down.sql` files in reverse order.
6. `make start` - Starts the app.
7. `make stop` - Stops the app.
8. `make build` - Builds the app into the `bin/infiniti` binary.
9. `make install` - Installs all uninstalled dependencies found in the app.
10. `make install PKG=<dependency>` - Installs a specific dependency. (e.g. `make install PKG=github.com/gorilla/mux`)
11. `make tidy` - Equivalent to `go mod tidy`, but executed directly inside the app's container.
12. `make shell` - Opens the interactive shell.