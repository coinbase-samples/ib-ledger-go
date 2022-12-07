### Amplify Ledger

Package providing a double-entry accounting ledger for a Coinbase Prime Introducing Brokers integration. It is built using Golang 1.19, 
and runs as a GRPC application server with a PostgreSQL Database.

For more information about Coinbase Prime: https://www.coinbase.com/prime

### Warning
This is a sample reference implementation, and as such is not built to be fully production-ready. 
Do not directly use this code in a customer facing application without subjecting it to significant load testing and a security review.

### Local Environment Setup

Required Installations:
* [Golang-Migrate](https://github.com/golang-migrate/migrate) - Migration library used to stand up the database
This is used in the Makefile commands
```
brew install golang-migrate
```
* [Docker](https://docs.docker.com/get-docker/) - Containers are used to run the Postgres Database locally. 
Installation instructions unnecessary as you probably have this already installed. If not, follow the link and download Docker.

## Running Ledger
To start up the local server, run the following:
```
make start-local
```

This will deploy the docker container for the database and spin up the server running by default at port 8445.

To run the integration tests:
```
make integ-test
```

This will run tests located under the test/integration path.
