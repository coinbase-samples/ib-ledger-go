.PHONY: dbinit

run:
	go run .

dbinit:
	go run cmd/dba/main.go 

compile:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protos/ledger/*.proto

complete_test:
	go run examples/complete/complete_transaction.go
