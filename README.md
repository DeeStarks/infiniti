# INFINITI (in progress)

Infiniti is a simple bank application system, designed on the basis of the Hexagonal Architecture.

### Setup

NB: Make sure to have docker installed and running. For docker installation, see [here](https://docs.docker.com/get-docker/).

- `git clone https://gihub.com/deestarks/infiniti.git`
- `cd infiniti`
- Create a `.env` file in the root directory of the project.
- Copy the content of the `.env.example` file to the `.env` file. Change the values of the variables according to your choice or leave them as they are.
- `make migrate-up` - Migrates the database schema.
- `make start` - To start the application.


### Commands
3. `make create-migrations` - Creates the migrations files for up and down into `adapters/framework/db/migrations`. Default name is a timestamp in format `YYYYMMDDHHMMSS`. To specify a different name, pass it as an argument to the command (e.g. `make create-migrations MIGRATION_NAME=example_name`).
4. `make migrate-up` - Runs all `*.up.sql` migration files sequentially.
5. `make migrate-down` - Same as `make migrate-up` but runs `*.down.sql` files in reverse order.
6. `make start` - Starts the app.
7. `make stop` - Stops the app.
8. `make build` - Builds the app into the `bin/infiniti` binary.
9. `make install` - Installs all uninstalled dependencies found in the app.
10. `make install PKG=<dependency>` - Installs a specific dependency. (e.g. `make install PKG=github.com/gorilla/mux`)r
11. `make tidy` - Equivalent to `go mod tidy`, but executed directly inside the app's container.
12. `make shell` - Opens an interactive shell inside the app's container.
13. `make test` - Recursively run all tests OR pass a directory as an argument to run the tests in a directory (e.g. `make test TESTDIR=github.com/deestarks/infiniti/adapters/framework/db` to run the tests in the `db` directory).
14. `make db-shell` - Opens an interactive shell for the database.