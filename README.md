# INFINITI (in progress)

Infiniti is a core banking system, designed on the basis of the Hexagonal Architecture.

API documentation can be found [here](https://documenter.getpostman.com/view/14444131/UVyswvHJ).

### Getting Started

NB: Make sure to have docker and docker-compose installed and running. For docker installation, see [here](https://docs.docker.com/get-docker/).

- `git clone https://gihub.com/deestarks/infiniti.git`
- `cd infiniti`
- Create a `.env` file in the root directory of the project and copy the content of the `.env.example` file into it. You can do this from the command line with `cp .env.example .env` in the root directory of the project.
- `make startd` - To start the application in the background.
- Generate a key using `make keygen`, copy the key and set it to the `RESTAPI_SECRET` environment variable in the `.env` file created. ***Note: The key is important as it is used to sign JWT tokens.***
- There's a job to update currency rates in the database every day at 6AM. To make it work, go to [Exchange Rates API](https://www.exchangerate-api.com/) to generate a free API key. Copy your API key, then use it to set the `EXCHANGE_RATE_API_KEY` environment variable in the `.env` file created. ***Note: This key is not important. If it's not set, the job will be skipped.***
- `make migrate-up` - To migrate database schema.
- `make cli-createadmin` - To create an admin user.

🚀🚀 Did all that without any errors? Great! The application will be available at `http://localhost:8000`. `curl http://localhost:8000/api/v1/` to see a welcome message.


### Commands
1. `make create-migrations` - Creates the migrations files for up and down into `adapters/framework/db/migrations`. Default name is a timestamp in format `YYYYMMDDHHMMSS`. To specify a different name, pass it as an argument to the command (e.g. `make create-migrations MIGRATION_NAME=example_name`).
2. `make migrate-up` - Runs all `*.up.sql` migration files sequentially.
3. `make migrate-down` - Same as `make migrate-up` but runs `*.down.sql` files in reverse order. To force migrate-down to a specific version - `make migrate-down MIGRATE_VERSION=<version>`.
4. `make startd` - Starts the application in the background.
5. `make start` - Starts the application in the foreground to view logs live.
6. `make stop` - Stops the app.
7. `make logs` - Shows the app logs.
8. `make build` - Builds the app into the `bin/infiniti` binary.
9. `make install` - Installs all uninstalled dependencies found in the app.
10. `make install PKG=<dependency>` - Installs a specific dependency. (e.g. `make install PKG=github.com/gorilla/mux`)
11. `make tidy` - Equivalent to `go mod tidy`, but executed directly inside the app's container.
12. `make shell` - Opens an interactive shell inside the app's container.
13. `make test` - Recursively run all tests OR pass a directory as an argument to run the tests in a directory (e.g. `make test TESTDIR=github.com/deestarks/infiniti/adapters/framework/db` to run the tests in the `db` directory).
14. `make db-shell` - Opens an interactive shell for the database.
15. `make keygen` - Generates a random 32-character string.
16. `make cli-createadmin` - To create an admin user from the command line.


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
│   ├── framework
│   │   └── db                      # Contains all database components (e.g. models, migrations, etc.)
│   │       ├── migrations          # Contains all migration files.
│   │       └── constants           # Contains reusable constants.
│   └── jobs                        # Job queuing layer.
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