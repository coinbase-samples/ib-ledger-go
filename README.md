# IB-Ledger-Go

Package providing a double-entry accounting ledger for a Coinbase Prime Introducing Brokers integration. It is built using Golang 1.19, 
and runs as a gRPC application server with [Amazon Quantum Ledger Database](https://aws.amazon.com/qldb/) (QLDB) and a PostgreSQL Database.

For more information about [Coinbase Prime](https://www.coinbase.com/prime)

## Warning
This is a sample reference implementation, and as such is not built to be fully production-ready. 
Do not directly use this code in a customer facing application without subjecting it to significant load testing and a security review.

## Required Installations

* [Golang-Migrate](https://github.com/golang-migrate/migrate) - Migration library used to stand up the database
* [Docker](https://docs.docker.com/get-docker/) - Containers are used to run the Postgres Database locally. 

## Running Ledger
### AWS Setup
This app requires a test Ledger for QLDB to be created in an AWS account.

The default profile is `sa-infra`, and the default name for the Ledger is `LedgerTest`

Your test account needs to have the following:
1. QLDB Ledger
2. Kinesis Stream
3. SQS Queue
4. QLDB Stream connecting QLDB to Kinesis
5. Eventbridge Pipe connecting Kinesis to SQS

To Setup QLDB Execute:
```
go run cmd/dba/initialize/main.go
```

This will set up the Ledger with two tables (Account and Ledger) and initialize the 
proper fee accounts.

### Start Application
To start up the local server, run the following:
```
make start-local
```

This does the following:
* Runs docker-compose up
* Runs migrate up
* Executes the test/sql directory file to insert testing data
* Starts the application server - default port is 8445

### Integration Test
Integration tests are stored under test/integration. These are currently configured to only run locally. Cleanup is still being implemented, so the 
local environment will need to be torn down and respun up in between test runs.

Before the first execution of the integration tests run:
```
go run test/setup/setup.go
```
This will setup the test accounts so that you can run the tests.

To execute the tests:
```
make integ-test
```

## License
This library is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file.

