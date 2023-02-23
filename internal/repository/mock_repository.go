/**
 * Copyright 2022 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package repository

import (
	"context"

	ledgererr "github.com/coinbase-samples/ib-ledger-go/internal/errors"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"google.golang.org/grpc/codes"

	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
)

type MockRepository struct {
}

var (
	CompleteTransactionUuidNotFound = "4AC6E407-1D8E-4339-BA1C-862ACC58AC5E"
	FailTransactionUuidNotFound     = "20032259-738B-40A7-AAD7-306B69AF88D4"
	CancelTransactionUuidNotFound   = "E6096F2D-C706-42B6-B0E5-D7DD644ED079"
)

func (m *MockRepository) InitializeAccount(ctx context.Context, request *api.InitializeAccountRequest) (*model.InitializeAccountResult, error) {
	return &model.InitializeAccountResult{}, nil
}

func (m *MockRepository) CreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (*model.CreateTransactionResult, error) {
	return &model.CreateTransactionResult{}, nil
}

func (m *MockRepository) PartialReleaseHold(ctx context.Context, request *api.PostFillRequest) (*model.TransactionResult, error) {
	return &model.TransactionResult{}, nil
}

func (m *MockRepository) CompleteTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	if request.OrderId == CompleteTransactionUuidNotFound {
		return nil, ledgererr.New(codes.NotFound, "LGR404")
	}

	return &model.TransactionResult{}, nil
}

func (m *MockRepository) FailTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	if request.OrderId == FailTransactionUuidNotFound {
		return nil, ledgererr.New(codes.NotFound, "LGR404")
	}

	return &model.TransactionResult{}, nil
}

func (m *MockRepository) CancelTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	if request.OrderId == CancelTransactionUuidNotFound {
		return nil, ledgererr.New(codes.NotFound, "LGR404")
	}

	return &model.TransactionResult{}, nil
}

func (m *MockRepository) GetAllAccountsAndMostRecentBalances(ctx context.Context, userId string) ([]*model.GetAccountResult, error) {
	return []*model.GetAccountResult{{}}, nil

}
