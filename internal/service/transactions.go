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
	"time"

	"github.com/amzn/ion-go/ion"
	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/coinbase-samples/ib-ledger-go/internal/qldb"
	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateTransaction(
	ctx context.Context,
	req *api.CreateTransactionRequest,
) (*api.CreateTransactionResponse, error) {
	l := ctxlogrus.Extract(ctx)

	if err := req.ValidateAll(); err != nil {
		l.Debugf("invalid create transaction request: %v", req)
		return nil, handleValidationError(err)
	}

	var orderAmount *big.Int
	if err := utils.SetString(orderAmount, req.TotalAmount); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"failed to convert TotalAmount to int - transaction: %s - err: %w",
			req.OrderId,
			err,
		)
	}

	feeAmount := big.NewInt(0)
	if req.FeeAmount != nil {
		if err := utils.SetString(
			feeAmount,
			req.FeeAmount.Value,
		); err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"failed to convert TotalAmount to int - transaction: %s - err: %w",
				req.OrderId,
				err,
			)
		}
	}

	var holdAmount *big.Int
	holdAmount.Add(orderAmount, feeAmount)

	senderUserId := strings.ToLower(req.Sender.UserId)
	senderId := model.GenerateAccountId(senderUserId, req.Sender.Currency)
	sender, err := qldb.GetAccount(ctx, senderId)

	if err != nil {
		return nil, handleTransactionErrors(err)
	}

	receiverUserId := strings.ToLower(req.Receiver.UserId)
	receiverCurrency := strings.ToUpper(req.Receiver.Currency)
	receiverId := model.GenerateAccountId(receiverUserId, receiverCurrency)
	receiver, err := qldb.GetAccount(ctx, receiverId)
	if err != nil {
		if _, ok := err.(*qldb.AccountNotFoundError); !ok {
			return nil, handleTransactionErrors(err)
		}
		if err := qldb.CreateAccountAndInitializeBalance(
			ctx,
			receiverUserId,
			receiverCurrency,
			big.NewInt(0)); err != nil {
			return nil, handleTransactionErrors(err)
		}
		receiver, err = qldb.GetAccount(ctx, receiverId)
		if err != nil {
			return nil, handleTransactionErrors(err)
		}
	}

	h := &model.QldbHold{
		HoldUUID:   uuid.New().String(),
		AccountId:  senderId,
		Amount:     ion.NewDecimal(holdAmount, 0, false),
		Released:   false,
		ReleasedAt: time.Time{},
	}

	ttype, ok := utils.GetStringFromTransactionType(req.TransactionType)
	if !ok {
		return nil, handleValidationError(
			fmt.Errorf("bad transaction type: %s", ttype))
	}

	venueOrderId := strings.ToLower(req.OrderId)
	transactionId := model.GenerateTransactionId(venueOrderId)
	l.Debugf("creating transaction with id: %s", transactionId)

	input := &model.QldbTransaction{
		Id:              transactionId,
		VenueOrderId:    venueOrderId,
		Sender:          sender.GetCoreAccount(),
		Receiver:        receiver.GetCoreAccount(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Hold:            h,
		TransactionType: ttype,
		Status:          config.TransactionStatusPending,
	}

	if err := qldb.CreateTransactionAndPlaceHold(
		ctx, input, holdAmount); err != nil {
		l.Debugf("unable to create transaction with input: %v", input)
		return nil, handleTransactionErrors(err)
	}

	return &api.CreateTransactionResponse{}, nil
}

func (s *Service) PostFill(
	ctx context.Context,
	req *api.PostFillRequest,
) (*api.PostFillResponse, error) {
	l := ctxlogrus.Extract(ctx)

	if err := req.ValidateAll(); err != nil {
		l.Debugf("invalid post fill request: %v", req)
		return nil, handleValidationError(err)
	}

	if err := qldb.FillOrder(req); err != nil {
		return nil, fmt.Errorf("ib-ledger-go: %w", err)
	}
	return &api.PostFillResponse{}, nil
}

func (s *Service) FinalizeTransaction(
	ctx context.Context,
	req *api.FinalizeTransactionRequest,
) (*api.FinalizeTransactionResponse, error) {
	l := ctxlogrus.Extract(ctx)

	if err := req.ValidateAll(); err != nil {
		l.Debugf("invalid finalize transaction request: %v", req)
		return nil, handleValidationError(err)
	}

	venueOrderId := strings.ToLower(req.OrderId)
	switch req.FinalizedStatus {
	case api.TransactionStatus_TRANSACTION_STATUS_COMPLETE:
		if err := qldb.FinalizeTransactionAndReleaseHold(
			ctx,
			venueOrderId,
			config.TransactionStatusFilled,
		); err != nil {
			return nil, handleTransactionErrors(err)
		}
	case api.TransactionStatus_TRANSACTION_STATUS_FAILED:
		if err := qldb.FinalizeTransactionAndReleaseHold(
			ctx,
			venueOrderId,
			config.TransactionStatusFailed,
		); err != nil {
			return nil, handleTransactionErrors(err)
		}
	case api.TransactionStatus_TRANSACTION_STATUS_CANCELED:
		if err := qldb.FinalizeTransactionAndReleaseHold(
			ctx,
			venueOrderId,
			config.TransactionStatusCanceled,
		); err != nil {
			return nil, handleTransactionErrors(err)
		}
	default:
		return nil, handleValidationError(
			fmt.Errorf(
				"unknown status: %s",
				req.FinalizedStatus,
			),
		)
	}

	return &api.FinalizeTransactionResponse{}, nil
}

func handleTransactionErrors(err error) error {
	wrapErr := fmt.Errorf("ib-ledger-go: %w", err)
	switch err.(type) {
	case *qldb.AccountNotFoundError:
		return status.Error(codes.NotFound, wrapErr.Error())
	case *qldb.InsufficientBalanceError:
		return status.Error(codes.InvalidArgument, wrapErr.Error())
	default:
		return status.Error(codes.Internal, wrapErr.Error())
	}
}

func handleValidationError(err error) error {
	return status.Error(
		codes.InvalidArgument,
		fmt.Sprintf("ib-ledger-go: %s", err.Error()),
	)
}
