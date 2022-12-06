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

package service

import (
	"context"
	"fmt"
	"log"

	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Service) InitializeAccount(ctx context.Context, req *api.InitializeAccountRequest) (*api.InitializeAccountResponse, error) {
	result, err := s.Repository.InitializeAccount(ctx, req)
	if err != nil {
		log.Printf("unable to initialize account %v", err)
		return nil, err
	}

	response := &api.InitializeAccountResponse{
		Account: &api.Account{
			Id:          &wrapperspb.StringValue{Value: result.Id.String()},
			PortfolioId: &wrapperspb.StringValue{Value: result.PortfolioId.String()},
			UserId:      result.UserId.String(),
			Currency:    result.Currency,
			CreatedAt:   timestamppb.New(result.CreatedAt),
		},
		AccountBalance: &api.AccountBalance{
			Balance:   fmt.Sprint(result.Balance),
			Hold:      fmt.Sprint(result.Hold),
			Available: fmt.Sprint(result.Available),
		},
	}
	return response, nil
}

func (s *Service) GetAccounts(ctx context.Context, req *api.GetAccountsRequest) (*api.GetAccountsResponse, error) {
	l := ctxlogrus.Extract(ctx)

	results, err := s.Repository.GetAllAccountsAndMostRecentBalances(ctx, req.UserId)
	if err != nil {
		l.Errorf("unable to get accounts and balances for user: %v", err)
		return nil, err
	}

	var outputResults []*api.AccountAndBalance

	for _, r := range results {
		outputResults = append(outputResults, &api.AccountAndBalance{
			AccountId: r.AccountId.String(),
			Currency:  r.Currency,
			Balance:   fmt.Sprint(r.Balance),
			Hold:      fmt.Sprint(r.Hold),
			Available: fmt.Sprint(r.Available),
			BalanceAt: timestamppb.New(r.CreatedAt),
		})
	}

	return &api.GetAccountsResponse{Accounts: outputResults}, nil
}
