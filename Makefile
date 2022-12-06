.PHONY: dbinit docker-build compile docker-build-local

REGION ?= us-east-1
PROFILE ?= sa-infra
ENV_NAME ?= dev
ACCOUNT_ID := $(shell aws sts get-caller-identity --profile $(PROFILE) --query 'Account' --output text)

.PHONY: run
run:
	go run cmd/server/main.go

compile:
	buf build && buf generate

.PHONY: complete-test
complete-test:
	go run examples/complete/complete_transaction.go

docker-build:
	@docker build --platform linux/amd64 --build-arg REGION=$(REGION) --build-arg ENV_NAME=$(ENV_NAME) --build-arg ACCOUNT_ID=$(ACCOUNT_ID) .

docker-build-local:
	@docker build --tag ib-ledger-go:local --build-arg REGION=$(REGION) --build-arg ENV_NAME=local --build-arg ACCOUNT_ID=$(ACCOUNT_ID) .

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
    && go run cmd/server/main.go

