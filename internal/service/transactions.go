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

	ledgererr "github.com/coinbase-samples/ib-ledger-go/internal/errors"
	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) CreateTransaction(ctx context.Context, req *api.CreateTransactionRequest) (*api.CreateTransactionResponse, error) {
	l := ctxlogrus.Extract(ctx)

	if err := req.ValidateAll(); err != nil {
		l.Debugf("invalid create transaction request: %v", req)
		return nil, fmt.Errorf("ib-ledger-go: %w", err)
	}

	result, err := s.Repository.CreateTransaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ib-ledger-go: %w", err)
	}
	transactionType, ok := utils.GetTransactionTypeFromString(result.TransactionType)
	if !ok {
		return nil, ledgererr.New(codes.InvalidArgument, fmt.Sprintf("ib-ledger-go: bad request: transaction type not supported: %v", result.TransactionType))
	}
	return &api.CreateTransactionResponse{
		Transaction: &api.Transaction{
			Id:                result.Id.String(),
			SenderId:          result.SenderId.String(),
			ReceiverId:        result.ReceiverId.String(),
			CreatedAt:         timestamppb.New(result.CreatedAt),
			TransactionType:   transactionType,
			RequestId:         result.RequestId.String(),
			TransactionStatus: api.TransactionStatus_TRANSACTION_STATUS_PENDING,
		},
	}, nil
}

func (s *Service) PartialReleaseHold(ctx context.Context, req *api.PartialReleaseHoldRequest) (*api.PartialReleaseHoldResponse, error) {
	l := ctxlogrus.Extract(ctx)

	if err := req.ValidateAll(); err != nil {
		l.Debugf("invalid partial release hold request: %v", req)
		return nil, fmt.Errorf("ib-ledger-go: %w", err)
	}

	_, err := s.Repository.PartialReleaseHold(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ib-ledger-go: %w", err)
	}
	return &api.PartialReleaseHoldResponse{
		Successful: true,
	}, nil
}

func (s *Service) FinalizeTransaction(ctx context.Context, req *api.FinalizeTransactionRequest) (*api.FinalizeTransactionResponse, error) {
	l := ctxlogrus.Extract(ctx)

	if err := req.ValidateAll(); err != nil {
		l.Debugf("invalid finalize transaction request: %v", req)
		return nil, fmt.Errorf("ib-ledger-go: %w", err)
	}

	switch req.FinalizedStatus {
	case api.TransactionStatus_TRANSACTION_STATUS_COMPLETE:
		_, err := s.Repository.CompleteTransaction(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("ib-ledger-go: %w", err)
		}
	case api.TransactionStatus_TRANSACTION_STATUS_FAILED:
		_, err := s.Repository.FailTransaction(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("ib-ledger-go: %w", err)
		}
	case api.TransactionStatus_TRANSACTION_STATUS_CANCELED:
		_, err := s.Repository.CancelTransaction(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("ib-ledger-go: %w", err)
		}
	default:
		return nil, ledgererr.New(codes.InvalidArgument, fmt.Sprintf("ib-ledger-go: finalize transaction: unable to finalize pending transaction - transaction: %v", req.OrderId))
	}

	return &api.FinalizeTransactionResponse{
		Successful: true,
	}, nil
}
