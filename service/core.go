package service

import (
	api "LedgerApp/protos/ledger"
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) InitializeAccount(ctx context.Context, req *api.InitializeAccountRequest) (*api.InitializeAccountResponse, error) {
	result, err := s.PostgresHandler.InitializeAccount(ctx, req)
	if err != nil {
		log.Printf("unable to initialize account %v", err)
		return nil, err
	}
	response := &api.InitializeAccountResponse{
		Account: &api.Account{
			Id:          result.Id.String(),
			PortfolioId: result.PortfolioId.String(),
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
	return nil, nil
}

func (s *Service) GetBalance(ctx context.Context, req *api.GetBalanceRequest) (*api.GetBalanceResponse, error) {
	return nil, nil
}

func (s *Service) GetBalances(ctx context.Context, req *api.GetBalancesRequest) (*api.GetBalancesResponse, error) {
	return nil, nil
}
