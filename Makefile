.PHONY: docker-build compile docker-build-local docker-setup start start-full restart-full integ-test

REGION ?= us-east-1
PROFILE ?= cb-stp-sa-dev
ENV_NAME ?= dev
ACCOUNT_ID := $(shell aws sts get-caller-identity --profile $(PROFILE) --query 'Account' --output text)

compile:
	@buf build 
	@buf generate

docker-build:
	@docker build --platform linux/amd64 --build-arg REGION=$(REGION) --build-arg ENV_NAME=$(ENV_NAME) --build-arg ACCOUNT_ID=$(ACCOUNT_ID) .

docker-build-local:
	@docker build --tag ib-ledger-go:local --build-arg REGION=$(REGION) --build-arg ENV_NAME=local --build-arg ACCOUNT_ID=$(ACCOUNT_ID) .

docker-setup:
	@docker-compose up -d
	@sleep 5
	@migrate -source file://dba/migrations/ -database 'postgres://postgres:postgres@localhost:5432/ledger?sslmode=disable' up
	@docker exec -i ledger_db psql -U postgres -d ledger < test/sql/initialize_test_accounts.sql

docker-shutdown:
	@docker-compose down

start:
	@go run cmd/server/*.go

start-full: docker-setup start

restart-full: docker-shutdown docker-setup start

integ-test:
	@go test ./test/integration

unit-test:
	@go test ./internal/...
