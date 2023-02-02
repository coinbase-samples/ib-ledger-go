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

package test

import (
	"context"
	"testing"

	"github.com/amzn/ion-go/ion"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/google/uuid"

	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestSuccessfulCompleteTransaction(t *testing.T) {
	ctx := context.Background()

	ledgerClient := newLedgerServiceClient(ctx, t)

	userId := uuid.NewString()
	setupAccounts(t, ctx, ledgerClient, userId)

	orderId := uuid.NewString()
	createTransactionRequest := &ledger.CreateTransactionRequest{
		OrderId: orderId,
		Sender: &ledger.Account{
			UserId:   userId,
			Currency: "USD",
		},
		Receiver: &ledger.Account{
			UserId:   userId,
			Currency: "ETH",
		},
		TotalAmount:     "10000",
		FeeAmount:       &wrapperspb.StringValue{Value: "100"},
		TransactionType: ledger.TransactionType_TRANSACTION_TYPE_TRANSFER,
		RequestId:       &wrapperspb.StringValue{Value: "85688730-9A34-4F9C-8475-7521B957F164"},
	}

	createTransaction(ledgerClient, ctx, t, createTransactionRequest)

	usdAccount := &model.QldbAccount{
		Id:        model.GenerateAccountId(userId, "USD"),
		UserId:    userId,
		Currency:  "USD",
		Balance:   ion.MustParseDecimal("100000"),
		Hold:      ion.MustParseDecimal("10100"),
		Available: ion.MustParseDecimal("89900"),
	}

	ethAccount := &model.QldbAccount{
		Id:        model.GenerateAccountId(userId, "ETH"),
		UserId:    userId,
		Currency:  "ETH",
		Balance:   ion.MustParseDecimal("0"),
		Hold:      ion.MustParseDecimal("0"),
		Available: ion.MustParseDecimal("0"),
	}

	getAccountsAndConfirmBalances(t, ctx, ledgerClient, userId, usdAccount, ethAccount)

	_, err := ledgerClient.PostFill(ctx, &ledger.PostFillRequest{
		OrderId:        orderId,
		FillId:         "11A40B5B-74AA-4CBD-A04B-AADC4F0487E8",
		FilledValue:    "1000",
		FilledQuantity: "1",
	})

	if err != nil {
		t.Fatalf("unable to post fill: %v", err)
	}

	usdAccount.Balance = ion.MustParseDecimal("99000")
	usdAccount.Hold = ion.MustParseDecimal("9100")
	usdAccount.Available = ion.MustParseDecimal("89900")
	ethAccount.Balance = ion.MustParseDecimal("1")
	ethAccount.Hold = ion.MustParseDecimal("0")
	ethAccount.Available = ion.MustParseDecimal("1")
	getAccountsAndConfirmBalances(t, ctx, ledgerClient, userId, usdAccount, ethAccount)

	_, err = ledgerClient.PostFill(ctx, &ledger.PostFillRequest{
		OrderId:         orderId,
		FillId:          uuid.NewString(),
		FilledValue:     "9000",
		FilledQuantity:  "9",
		VenueFeeAmount:  &wrapperspb.StringValue{Value: "50"},
		RetailFeeAmount: &wrapperspb.StringValue{Value: "50"},
	})

	if err != nil {
		t.Fatalf("unable to post fill: %v", err)
	}

	usdAccount.Balance = ion.MustParseDecimal("89900")
	usdAccount.Hold = ion.MustParseDecimal("0")
	usdAccount.Available = ion.MustParseDecimal("89900")
	ethAccount.Balance = ion.MustParseDecimal("10")
	ethAccount.Hold = ion.MustParseDecimal("0")
	ethAccount.Available = ion.MustParseDecimal("10")
	getAccountsAndConfirmBalances(t, ctx, ledgerClient, userId, usdAccount, ethAccount)

	if _, err := ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         orderId,
		FinalizedStatus: ledger.TransactionStatus_TRANSACTION_STATUS_COMPLETE,
	}); err != nil {
		t.Fatalf("unable to finalize transaction: %v", err)
	}

	getAccountsAndConfirmBalances(t, ctx, ledgerClient, userId, usdAccount, ethAccount)

	orderId = uuid.NewString()
	createTransactionRequest = &ledger.CreateTransactionRequest{
		OrderId: orderId,
		Sender: &ledger.Account{
			UserId:   userId,
			Currency: "ETH",
		},
		Receiver: &ledger.Account{
			UserId:   userId,
			Currency: "USD",
		},
		TotalAmount:     "10",
		TransactionType: ledger.TransactionType_TRANSACTION_TYPE_TRANSFER,
		RequestId:       &wrapperspb.StringValue{Value: "85688730-9A34-4F9C-8475-7521B957F164"},
	}

	createTransaction(ledgerClient, ctx, t, createTransactionRequest)

	usdAccount.Balance = ion.MustParseDecimal("89900")
	usdAccount.Hold = ion.MustParseDecimal("0")
	usdAccount.Available = ion.MustParseDecimal("89900")
	ethAccount.Balance = ion.MustParseDecimal("10")
	ethAccount.Hold = ion.MustParseDecimal("10")
	ethAccount.Available = ion.MustParseDecimal("0")
	getAccountsAndConfirmBalances(t, ctx, ledgerClient, userId, usdAccount, ethAccount)

	_, err = ledgerClient.PostFill(ctx, &ledger.PostFillRequest{
		OrderId:        orderId,
		FillId:         uuid.NewString(),
		FilledValue:    "10",
		FilledQuantity: "10100",
	})

	if err != nil {
		t.Fatalf("unable to post fill: %v", err)
	}

	usdAccount.Balance = ion.MustParseDecimal("100000")
	usdAccount.Hold = ion.MustParseDecimal("0")
	usdAccount.Available = ion.MustParseDecimal("100000")
	ethAccount.Balance = ion.MustParseDecimal("0")
	ethAccount.Hold = ion.MustParseDecimal("0")
	ethAccount.Available = ion.MustParseDecimal("0")
	getAccountsAndConfirmBalances(t, ctx, ledgerClient, userId, usdAccount, ethAccount)

	if _, err := ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         orderId,
		FinalizedStatus: ledger.TransactionStatus_TRANSACTION_STATUS_COMPLETE,
	}); err != nil {
		t.Fatalf("unable to finalize transaction: %v", err)
	}

	getAccountsAndConfirmBalances(t, ctx, ledgerClient, userId, usdAccount, ethAccount)
}
