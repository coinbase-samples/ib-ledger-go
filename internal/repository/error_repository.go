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
	"errors"

	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
)

type ErrorRepository struct {
    underlying Repository
}

func NewErrorRepository(underlying Repository) Repository {
    return &ErrorRepository{
        underlying: underlying,
    }
}

func (r *ErrorRepository) handleCommonErrors(err error) error {
    var pgErr *pgconn.PgError

    if errors.As(err, &pgErr) {
        msg := pgErr.Message
        log.Errorf("postgres error - code: %v - message: %v", pgErr.SQLState(), msg)
        if msg == "insufficient available balance" {
            return errors.New("insufficient balance for transaction")
        } else if msg == "sender account missing" || msg == "receiver account missing" {
            return errors.New("account not found")
        } else {
            return errors.New("postgres internal failure")
        }
    }
    return err
}

func (r *ErrorRepository) InitializeAccount(ctx context.Context, request *api.InitializeAccountRequest) (*model.InitializeAccountResult, error) {
    resp, err := r.underlying.InitializeAccount(ctx, request)
    if err != nil {
        log.Errorf("error InitializeAccount - currency: %v - userId: %v", request.Currency, request.UserId)
        return resp, r.handleCommonErrors(err)
    }
    return resp, nil 
}

func (r *ErrorRepository) CreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (*model.CreateTransactionResult, error) {
    resp, err := r.underlying.CreateTransaction(ctx, request)
    if err != nil {
        log.Errorf("error CreateTransaction - orderId: %v", request.OrderId)
        return resp, r.handleCommonErrors(err)
    }
    return resp, nil 
}

func (r *ErrorRepository) PartialReleaseHold(ctx context.Context, request *api.PartialReleaseHoldRequest) (*model.TransactionResult, error) {
    resp, err := r.underlying.PartialReleaseHold(ctx, request)
    if err != nil {
        log.Errorf("error PartialReleaseHold - orderId: %v", request.OrderId)
        return resp, r.handleCommonErrors(err)
    }
    return resp, nil 
}

func (r *ErrorRepository) CompleteTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
    resp, err := r.underlying.CompleteTransaction(ctx, request)
    if err != nil {
        log.Errorf("error CompleteTransaction - orderId: %v", request.OrderId)
        return resp, r.handleCommonErrors(err)
    }
    return resp, nil 
}

func (r *ErrorRepository) FailTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
    resp, err := r.underlying.FailTransaction(ctx, request)
    if err != nil {
        log.Errorf("error FailTransaction - orderId: %v", request.OrderId)
        return resp, r.handleCommonErrors(err)
    }
    return resp, nil 
}

func (r *ErrorRepository) CancelTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
    resp, err := r.underlying.CancelTransaction(ctx, request)
    if err != nil {
        log.Errorf("error CancelTransaction - orderId: %v", request.OrderId)
        return resp, r.handleCommonErrors(err)
    }
    return resp, nil 
}

func (r *ErrorRepository) GetAllAccountsAndMostRecentBalances(ctx context.Context, userId string) ([]*model.GetAccountResult, error) {
    resp, err := r.underlying.GetAllAccountsAndMostRecentBalances(ctx, userId)
    if err != nil {
        log.Errorf("error GetAllAccountsAndMostRecentBalances - req: %v", userId)
        return resp, r.handleCommonErrors(err)
    }
    return resp, nil 
}

