.PHONY: dbinit docker-build compile

run:
	go run .

compile:
	buf build && buf generate

.PHONY: complete-test
complete-test:
	go run examples/complete/complete_transaction.go

docker-build:
	@docker build --platform linux/amd64 --build-arg REGION=$(REGION) --build-arg ENV_NAME=$(ENV_NAME) --build-arg ACCOUNT_ID=$(ACCOUNT_ID) .

.PHONY: init-test
init-test:
	migrate -source file://dba/migrations/ -database 'postgres://postgres:postgres@localhost:5432/ledger?sslmode=disable' up \
	&& docker exec -i ledger_db psql -U postgres -d ledger < examples/sql/initialize_test_accounts.sql

.PHONY: migrate-down-up
migrate-down-up:
	migrate -source file://dba/migrations/ -database 'postgres://postgres:postgres@localhost:5432/ledger?sslmode=disable' down \
		&& migrate -source file://dba/migrations/ -database 'postgres://postgres:postgres@localhost:5432/ledger?sslmode=disable' up \
    	&& docker exec -i ledger_db psql -U postgres -d ledger < examples/sql/initialize_test_accounts.sql \
    	&& go run .

.PHONY: start-local
start-local:
	docker-compose up -d \
	&& sleep 5 \
	&& migrate -source file://dba/migrations/ -database 'postgres://postgres:postgres@localhost:5432/ledger?sslmode=disable' up \
    && docker exec -i ledger_db psql -U postgres -d ledger < examples/sql/initialize_test_accounts.sql \
    && go run .
