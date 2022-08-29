package service

import (
	api "LedgerApp/protos/ledger"
	"LedgerApp/utils"
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) CreateTransaction(ctx context.Context, req *api.CreateTransactionRequest) (*api.CreateTransactionResponse, error) {
	result, err := s.PostgresHandler.CreateTransaction(ctx, req)
	if err != nil {
		log.Printf("unable to initialize account %v", err)
		return nil, err
	}
	response := &api.CreateTransactionResponse{
		Transaction: &api.Transaction{
			Id:                result.Id.String(),
			SenderId:          result.SenderId.String(),
			ReceiverId:        result.ReceiverId.String(),
			CreatedAt:         timestamppb.New(result.CreatedAt),
			TransactionType:   utils.GetTransactionTypeFromString(result.TransactionType),
			RequestId:         result.RequestId.String(),
			TransactionStatus: api.TransactionStatus_PENDING,
		},
	}
	return response, nil
}

func (s *Service) PartialReleaseHold(ctx context.Context, req *api.PartialReleaseHoldRequest) (*api.PartialReleaseHoldResponse, error) {
	_, err := s.PostgresHandler.PartialReleaseHold(ctx, req)
	if err != nil {
		log.Printf("unable to partial release hold: %v", err)
		return nil, err
	}
	response := &api.PartialReleaseHoldResponse{
		Successful: true,
	}
	return response, nil
}

func (s *Service) FinalizeTransaction(ctx context.Context, req *api.FinalizeTransactionRequest) (*api.FinalizeTransactionResponse, error) {
	switch req.FinalizedStatus {
	case api.TransactionStatus_COMPLETE:
		_, err := s.PostgresHandler.CompleteTransaction(ctx, req)
		if err != nil {
			log.Printf("unable to complete transaction: %v", err)
			return nil, err
		}
	case api.TransactionStatus_FAILED:
		_, err := s.PostgresHandler.FailTransaction(ctx, req)
		if err != nil {
			log.Printf("unable to fail transaction: %v", err)
			return nil, err
		}
	case api.TransactionStatus_CANCELED:
		_, err := s.PostgresHandler.CancelTransaction(ctx, req)
		if err != nil {
			log.Printf("unable to cancel transaction: %v", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("finalize transaction: unsupported finalized transaction status %v", req.FinalizedStatus)
	}

	response := &api.FinalizeTransactionResponse{
		Successful: true,
	}
	return response, nil
}
