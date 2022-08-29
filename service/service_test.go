package service

import (
	"LedgerApp/model"
	"context"

	api "LedgerApp/protos/ledger"
)

type TestPostgres struct {
}

func (tp TestPostgres) InitializeAccount(ctx context.Context, request *api.InitializeAccountRequest) (*model.InitializeAccountResult, error) {
	return &model.InitializeAccountResult{}, nil
}
