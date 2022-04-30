# INFINITI (in progress)

Infiniti is a simple bank application system, designed on the basis of the Hexagonal Architecture.

API documentation can be found [here](https://documenter.getpostman.com/view/14444131/UVyswvHJ).

### Setup

NB: Make sure to have docker installed and running. For docker installation, see [here](https://docs.docker.com/get-docker/).

- `git clone https://gihub.com/deestarks/infiniti.git`
- `cd infiniti`
- Create a `.env` file in the root directory of the project.
- Copy the content of the `.env.example` file to the `.env` file. Change the values of the variables according to your choice or leave them as they are. ***(Important: the `RESTAPI_SECRET` variable is important. You can generate one using `make keygen` command)***
- `make start` - To start the application.
- `make migrate-up` - Migrates the database schema.
- `make cli-createadmin` - To create an admin user.

App will be available at `http://localhost:8000`. `curl http://localhost:8000/api/v1/` will return a welcome message.


### Commands
1. `make create-migrations` - Creates the migrations files for up and down into `adapters/framework/db/migrations`. Default name is a timestamp in format `YYYYMMDDHHMMSS`. To specify a different name, pass it as an argument to the command (e.g. `make create-migrations MIGRATION_NAME=example_name`).
2. `make migrate-up` - Runs all `*.up.sql` migration files sequentially.
3. `make migrate-down` - Same as `make migrate-up` but runs `*.down.sql` files in reverse order. To force migrate-down to a specific version - `make migrate-down MIGRATE_VERSION=<version>`.
4. `make start` - Starts the app.
5. `make stop` - Stops the app.
6. `make log` - Shows the app logs.
7. `make build` - Builds the app into the `bin/infiniti` binary.
8. `make install` - Installs all uninstalled dependencies found in the app.
9. `make install PKG=<dependency>` - Installs a specific dependency. (e.g. `make install PKG=github.com/gorilla/mux`)
10. `make tidy` - Equivalent to `go mod tidy`, but executed directly inside the app's container.
11. `make shell` - Opens an interactive shell inside the app's container.
12. `make test` - Recursively run all tests OR pass a directory as an argument to run the tests in a directory (e.g. `make test TESTDIR=github.com/deestarks/infiniti/adapters/framework/db` to run the tests in the `db` directory).
13. `make db-shell` - Opens an interactive shell for the database.
14. `make keygen` - Generates a random 32-character string.
15. `make cli-createadmin` - To create an admin user from the command line.


### Project structure
```
infiniti
├── adapters
│   ├── client
│   │    ├── restapi                # Contains all REST API app components (e.g. handlers, middlewares, etc.)
│   │    │   ├── constants          # Constants used accross the RESTAPI client.
│   │    │   ├── handlers           # Contains all handlers for the REST API including the routes.
│   │    │   └── middleware         # Contains all middlewares for the REST API.
│   │    └── cli                    # Contains all CLI app components (e.g. handlers, etc.)
│   │       └── handlers            # Contains all handlers for the CLI app.
│   └── framework
│       └── db                      # Contains all database components (e.g. models, migrations, etc.)
│           ├── migrations          # Contains all migration files.
│           └── constants           # Contains reusable constants.
├── application                     # Application domain.
│   ├── core                        # Contains business logic.
│   │   └── tests                   # Tests for the core domain.
│   └── services                    # Contains application services
├── bin                             # Contains the binary for the app.
│   └── infiniti
├── cmd                             # Application entry point.
├── config                          # Contains all configuration files.
├── utils                           # Contains utility functions.
└── vol                             # Dedicated folder for docker volumes.
```