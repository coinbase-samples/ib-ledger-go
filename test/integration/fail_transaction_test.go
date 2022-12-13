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

package test

import (
	"context"
	"testing"

	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func FailTransactionSucceedsTest(t *testing.T) {
	ctx := context.Background()

	ledgerClient := newLedgerServiceClient(ctx, t)

	orderId := "F2F83699-C5BD-405D-9ACE-1FAD2A3DEB44"
	createTransactionRequest := &ledger.CreateTransactionRequest{
		OrderId: orderId,
		Sender: &ledger.Account{
			UserId:   "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
			Currency: "USD",
		},
		Receiver: &ledger.Account{
			UserId:   "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
			Currency: "ETH",
		},
		TotalAmount:     "1000",
		TransactionType: ledger.TransactionType_TRANSACTION_TYPE_TRANSFER,
		RequestId:       &wrapperspb.StringValue{Value: "D6F45D7D-24E4-435A-ADF7-1BE9BDD081DE"},
	}

	createTransactionAndConfirmHolds(ledgerClient, ctx, t, createTransactionRequest)

	senderExpectedBalance := &ledger.AccountAndBalance{
		Currency:  "USD",
		Balance:   "100000",
		Hold:      "1000",
		Available: "99000",
	}
	receiverExpectedBalance := &ledger.AccountAndBalance{
		Currency:  "ETH",
		Balance:   "100000",
		Hold:      "0",
		Available: "100000",
	}

	getTransactionBalancesAndConfirmTheyAreAsExpected(ledgerClient, ctx, t, senderExpectedBalance, receiverExpectedBalance)

	partialReleaseHoldAndConfirmSuccessful(ledgerClient, ctx, t, &ledger.PartialReleaseHoldRequest{
		OrderId:        orderId,
		RequestId:      "4FC0EA00-1709-491B-BFD3-7EA85B9C7447",
		SenderAmount:   "1",
		ReceiverAmount: "10",
	})

	senderExpectedBalance = &ledger.AccountAndBalance{
		Currency:  "USD",
		Balance:   "99999",
		Hold:      "999",
		Available: "99000",
	}
	receiverExpectedBalance = &ledger.AccountAndBalance{
		Currency:  "ETH",
		Balance:   "100010",
		Hold:      "0",
		Available: "100010",
	}
	getTransactionBalancesAndConfirmTheyAreAsExpected(ledgerClient, ctx, t, senderExpectedBalance, receiverExpectedBalance)

	finalizeResult, err := ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         orderId,
		RequestId:       "28367E08-98A4-4614-9FA6-FCD830CCA9F7",
		FinalizedStatus: ledger.TransactionStatus_TRANSACTION_STATUS_FAILED,
	})

	if err != nil {
		t.Fatalf("unable to fail transaction: %v", err.Error())
	}

	assert.True(t, finalizeResult.Successful)

	senderExpectedBalance = &ledger.AccountAndBalance{
		Currency:  "USD",
		Balance:   "99999",
		Hold:      "0",
		Available: "99999",
	}
	receiverExpectedBalance = &ledger.AccountAndBalance{
		Currency:  "ETH",
		Balance:   "100010",
		Hold:      "0",
		Available: "100010",
	}
	getTransactionBalancesAndConfirmTheyAreAsExpected(ledgerClient, ctx, t, senderExpectedBalance, receiverExpectedBalance)
}
