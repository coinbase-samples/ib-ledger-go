/**
 * Copyright 2022-present Coinbase Global, Inc.
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

package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/coinbase-samples/ib-ledger-go/internal/qldb"
	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) InitializeAccount(
	ctx context.Context,
	req *api.InitializeAccountRequest,
) (*api.InitializeAccountResponse, error) {
	l := ctxlogrus.Extract(ctx)

	if err := req.ValidateAll(); err != nil {
		l.Debugf("invalid initialize account request: %v", req)
		return nil, handleValidationError(err)
	}

	var initialBalance *big.Int
	if req.InitialBalance == "" {
		initialBalance = big.NewInt(0)
	} else {
		var initialBalance *big.Int
		if err := utils.SetString(initialBalance, req.InitialBalance); err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"failed to convert TotalAmount to int - userId: %s - currency: %s - err: %w",
				req.UserId,
				req.Currency,
				err,
			)
		}
		l.Debugf("userId: %s, initial balance: %s", req.UserId, initialBalance.String())
	}
	err := qldb.CreateAccountAndInitializeBalance(
		ctx,
		strings.ToLower(req.UserId),
		strings.ToUpper(req.Currency),
		initialBalance,
	)
	if err != nil {
		return nil, handleTransactionErrors(err)
	}

	response := &api.InitializeAccountResponse{}
	return response, nil
}

func (s *Service) GetAccounts(
	ctx context.Context,
	req *api.GetAccountsRequest,
) (*api.GetAccountsResponse, error) {
	l := ctxlogrus.Extract(ctx)

	if err := req.ValidateAll(); err != nil {
		l.Debugf("invalid get accounts request: %v", req)
		return nil, handleValidationError(err)
	}

	results, err := qldb.GetUserAccounts(ctx, req.UserId)
	if err != nil {
		return nil, handleTransactionErrors(err)
	}

	var outputResults []*api.AccountAndBalance

	for _, r := range results {
		balance, err := utils.IonDecimalToBigInt(r.Balance)
		if err != nil {
			return nil, handleValidationError(
				fmt.Errorf(
					"bad balance in account balance result - id: %s - val: %s - %w",
					r.Id,
					r.Balance.String(),
					err,
				),
			)
		}

		hold, err := utils.IonDecimalToBigInt(r.Hold)
		if err != nil {
			return nil, handleValidationError(
				fmt.Errorf(
					"bad hold in account balance result - id: %s - val: %s - %w",
					r.Id,
					r.Hold.String(),
					err,
				),
			)
		}

		available, err := utils.IonDecimalToBigInt(r.Available)
		if err != nil {
			return nil, handleValidationError(
				fmt.Errorf(
					"bad available in account balance result - id: %s - val: %s - %w",
					r.Id,
					r.Available.String(),
					err,
				),
			)
		}

		time := r.UpdatedAt
		outputResults = append(
			outputResults,
			&api.AccountAndBalance{
				AccountId: r.Id,
				Currency:  r.Currency,
				Balance:   balance.String(),
				Hold:      hold.String(),
				Available: available.String(),
				BalanceAt: timestamppb.New(time),
			},
		)
	}

	return &api.GetAccountsResponse{Accounts: outputResults}, nil
}
