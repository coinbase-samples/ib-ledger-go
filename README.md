# Amplify Ledger

Package providing a double-entry accounting ledger for a Coinbase Prime Introducing Brokers integration. It is built using Golang 1.19, 
and runs as a GRPC application server with a PostgreSQL Database.

For more information about Coinbase Prime: https://www.coinbase.com/prime

## Warning
This is a sample reference implementation, and as such is not built to be fully production-ready. 
Do not directly use this code in a customer facing application without subjecting it to significant load testing and a security review.

## Required Installations

* [Golang-Migrate](https://github.com/golang-migrate/migrate) - Migration library used to stand up the database
* [Docker](https://docs.docker.com/get-docker/) - Containers are used to run the Postgres Database locally. 

## Running Ledger
### Start Application
To start up the local server, run the following:
```
make start-local
```

This does the following:
* runs docker-compose up
* runs migrate up
* executes the test/sql directory file to insert testing data
* starts the application server - default port is 8445

### Integration Test
Integration tests are stored under test/integration. These are currently configured to only run locally. Cleanup is still being implemented, so the 
local environment will need to be torn down and respun up in between test runs.

To execute the tests:
```
make integ-test
```

## License
This library is licensed under the Apache 2.0 License. See the LICENSE file.

