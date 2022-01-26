# INFINITI

Infiniti is a bank application for user accounts and transactions management. The system is designed on the basis of the Hexagonal Architecture.

### Setup

NB: Make sure to have docker installed and running. For docker installation, see [here](https://docs.docker.com/get-docker/).

- `git clone https://gihub.com/deestarks/infiniti.git`
- `cd infiniti`
- Create a `.env` file in the root directory of the project.
- Copy the content of the `.env.example` file to the `.env` file. Change the values of the variables according to your choice or leave them as they are.
- `make docker-build` - To build the docker image
- `make migrate-up` - Migrates the database schema. For more info about this command, see no. 4 in the [List of Commands](#commands) section.
- `make start` - To start the application.


### Commands
1. `make docker-build` - Builds the docker image, and automatically installs all dependencies and sets up the database.
2. `make db-setup` - Creates the database with values set in the `.env`. NB: You might not need to run this command because the `make docker-build` command will create the database automatically.
3. `make create-migrations` - Creates the migrations files for up and down into `adapters/framework/db/migrations`. Default name is a timestamp in format `YYYYMMDDHHMMSS`. To specify a different name, pass it as an argument to the command (e.g. `make create-migrations MIGRATION_NAME=example_name`).
4. `make migrate-up` - Runs all `*.up.sql` migration files sequentially.
5. `make migrate-down` - Same as `make migrate-up` but runs `*.down.sql` files in reverse order.
6. `make start` - Starts the app.
7. `make stop` - Stops the app.
8. `make build` - Builds the app into the `bin/infiniti` binary.
9. `make install` - Installs all uninstalled dependencies found in the app.
10. `make install PKG=<dependency>` - Installs a specific dependency. (e.g. `make install PKG=github.com/gorilla/mux`)
11. `make tidy` - Equivalent to `go mod tidy`, but executed directly inside the app's container.
12. `make shell` - Opens the interactive shell.
13. `make test` - Recursively run all tests OR pass a directory as an argument to run the tests in a directory (e.g. `make test TESTDIR=github.com/deestarks/infiniti/adapters/framework/db` to run the tests in the `db` directory).
14. `make db-shell` - Opens the interactive shell into the database.