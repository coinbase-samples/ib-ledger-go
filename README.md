### Amplify Ledger

This is the application code and SQL for the Amplify Ledger implementation.

### Local Environment Setup

Required Installations:
* [Golang-Migrate](https://github.com/golang-migrate/migrate) - Migration library used to stand up the database
```
brew install golang-migrate
```
* [Docker](https://docs.docker.com/get-docker/) - Containers are used to run the Postgres Database locally. 
Installation instructions unnecessary as you probably have this already installed. If not, follow the link and download Docker.

To spin up the container to run postgres:
`docker-compose up -d`

This will start up the database running at localhost:5432

To insert the database configuration and start the application layer run:

```
make start-local
```

This will insert the table schema and test data

To test the complete functionality, run the following command

```
make complete_test
```
