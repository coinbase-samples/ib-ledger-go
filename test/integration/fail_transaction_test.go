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
	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFailTransaction(t *testing.T) {
	ctx := context.Background()

	ledgerClient := newLedgerServiceClient(ctx, t)

	userId := uuid.New().String()
	setupAccounts(t, ctx, ledgerClient, userId)

	orderId := uuid.New().String()
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

	if _, err := ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         orderId,
		FinalizedStatus: ledger.TransactionStatus_TRANSACTION_STATUS_FAILED,
	}); err != nil {
		t.Fatalf("unable to finalize transaction: %v", err)
	}

	usdAccount.Balance = ion.MustParseDecimal("100000")
	usdAccount.Hold = ion.MustParseDecimal("0")
	usdAccount.Available = ion.MustParseDecimal("100000")
	ethAccount.Balance = ion.MustParseDecimal("0")
	ethAccount.Hold = ion.MustParseDecimal("0")
	ethAccount.Available = ion.MustParseDecimal("0")
	getAccountsAndConfirmBalances(t, ctx, ledgerClient, userId, usdAccount, ethAccount)

}
