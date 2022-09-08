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
	"LedgerApp/utils"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) CreateTransaction(ctx context.Context, req *api.CreateTransactionRequest) (*api.CreateTransactionResponse, error) {
	result, err := s.PostgresHandler.CreateTransaction(ctx, req)
	if err != nil {
		log.Errorf("unable to create transaction: %v", err)
		return nil, err
	}
	transactionType, ok := utils.GetTransactionTypeFromString(result.TransactionType)
	if !ok {
		return nil, fmt.Errorf("bad request: transaction type not supported: %v", result.TransactionType)
	}
	response := &api.CreateTransactionResponse{
		Transaction: &api.Transaction{
			Id:                result.Id.String(),
			SenderId:          result.SenderId.String(),
			ReceiverId:        result.ReceiverId.String(),
			CreatedAt:         timestamppb.New(result.CreatedAt),
			TransactionType:   transactionType,
			RequestId:         result.RequestId.String(),
			TransactionStatus: api.TransactionStatus_PENDING,
		},
	}
	return response, nil
}

func (s *Service) PartialReleaseHold(ctx context.Context, req *api.PartialReleaseHoldRequest) (*api.PartialReleaseHoldResponse, error) {
	_, err := s.PostgresHandler.PartialReleaseHold(ctx, req)
	if err != nil {
		log.Errorf("unable to partial release hold: %v", err)
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
			log.Errorf("unable to complete transaction: %v", err)
			return nil, err
		}
	case api.TransactionStatus_FAILED:
		_, err := s.PostgresHandler.FailTransaction(ctx, req)
		if err != nil {
			log.Errorf("unable to fail transaction: %v", err)
			return nil, err
		}
	case api.TransactionStatus_CANCELED:
		_, err := s.PostgresHandler.CancelTransaction(ctx, req)
		if err != nil {
			log.Errorf("unable to cancel transaction: %v", err)
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
