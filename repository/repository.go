package repository

import (
	"LedgerApp/model"
	api "LedgerApp/protos/ledger"
	"context"
)

type Repository interface {
	InitializeAccount(ctx context.Context, request *api.InitializeAccountRequest) (*model.InitializeAccountResult, error)
	CreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (*model.CreateTransactionResult, error)
	PartialReleaseHold(ctx context.Context, request *api.PartialReleaseHoldRequest) (*model.TransactionResult, error)
	CompleteTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error)
	FailTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error)
	CancelTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error)
}
