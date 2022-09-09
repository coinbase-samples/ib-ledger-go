### Amplify Ledger

This is the application code and SQL for the Amplify Ledger compute and database.

### Local Environment Setup

Environment Variables:
DB_HOSTNAME="localhost"
DB_PORT="5432"
DB_CREDENTIALS="{\"password\":\"postgres\",\"username\":\"postgres\"}"
ENV_NAME="local"

To spin up the database:
`docker-compose up -d`

This will start up the database running at localhost:5432

To insert the database configuration run the following command:

```
go run cmd/dba/main.go -m <PATH TO MIGRATIONS FOLDER> migrate up
```

This will create all of the tables and functions needed for the application.

To test the complete functionality, run the following command

```
make run complete_test
```
