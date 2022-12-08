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

func TestSuccessfulCompleteTransaction(t *testing.T) {
	ctx := context.Background()

	ledgerClient := newLedgerServiceClient(ctx, t)

	orderId := "10D0DF21-88E0-4076-8CD7-58F761E26F40"
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
		TransactionType: ledger.TransactionType_TRANSFER,
		RequestId:       &wrapperspb.StringValue{Value: "85688730-9A34-4F9C-8475-7521B957F164"},
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
		RequestId:      "11A40B5B-74AA-4CBD-A04B-AADC4F0487E8",
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

	partialReleaseHoldAndConfirmSuccessful(ledgerClient, ctx, t, &ledger.PartialReleaseHoldRequest{
		OrderId:        orderId,
		RequestId:      "7795D23A-D273-442B-B588-CD718E11E2E5",
		SenderAmount:   "999",
		ReceiverAmount: "10",
	})

	senderExpectedBalance = &ledger.AccountAndBalance{
		Currency:  "USD",
		Balance:   "99000",
		Hold:      "0",
		Available: "99000",
	}
	receiverExpectedBalance = &ledger.AccountAndBalance{
		Currency:  "ETH",
		Balance:   "100020",
		Hold:      "0",
		Available: "100020",
	}

	getTransactionBalancesAndConfirmTheyAreAsExpected(ledgerClient, ctx, t, senderExpectedBalance, receiverExpectedBalance)

	finalizeResponse, err := ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         orderId,
		RequestId:       "E2586CC8-2879-48B2-83E2-F1F5CF87E248",
		FinalizedStatus: ledger.TransactionStatus_COMPLETE,
	})

	if err != nil {
		t.Fatalf("unable to complete transaction")
	}

	assert.True(t, finalizeResponse.Successful)

	getTransactionBalancesAndConfirmTheyAreAsExpected(ledgerClient, ctx, t, senderExpectedBalance, receiverExpectedBalance)
}
