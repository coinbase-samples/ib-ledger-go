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
	api "LedgerApp/protos/ledger"
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Service) InitializeAccount(ctx context.Context, req *api.InitializeAccountRequest) (*api.InitializeAccountResponse, error) {
	result, err := s.PostgresHandler.InitializeAccount(ctx, req)
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

func (s *Service) GetAccount(ctx context.Context, req *api.GetAccountRequest) (*api.GetAccountResponse, error) {
	return nil, nil
}

func (s *Service) GetAccounts(ctx context.Context, req *api.GetAccountsRequest) (*api.GetAccountsResponse, error) {
	results, err := s.PostgresHandler.GetAllAccountsAndMostRecentBalances(ctx, req.UserId)
	if err != nil {
		log.Printf("unable to get accounts and balances for user: %v", err)
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

func (s *Service) GetBalance(ctx context.Context, req *api.GetBalanceRequest) (*api.GetBalanceResponse, error) {
	return nil, nil
}

func (s *Service) GetBalances(ctx context.Context, req *api.GetBalancesRequest) (*api.GetBalancesResponse, error) {
	return nil, nil
}
