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
)

func TestCompleteTransaction(t *testing.T) {
	ctx := context.Background()

	ledgerClient := newLedgerServiceClient(ctx, t)

	orderId := "10D0DF21-88E0-4076-8CD7-58F761E26F40"
	createTransactionRequestId := "85688730-9A34-4F9C-8475-7521B957F164"

	createTransactionAndConfirmHolds(ledgerClient, ctx, t, orderId, createTransactionRequestId)

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

	_, err := ledgerClient.PartialReleaseHold(ctx, &ledger.PartialReleaseHoldRequest{
		OrderId:        orderId,
		RequestId:      "11A40B5B-74AA-4CBD-A04B-AADC4F0487E8",
		SenderAmount:   "1",
		ReceiverAmount: "10",
	})

	if err != nil {
		t.Fatalf("unable to partial release hold")
	}

	output, err := ledgerClient.GetAccounts(ctx, &ledger.GetAccountsRequest{
		UserId: "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
	})

	if err != nil {
		t.Fatalf("complete transaction test: unable get accounts: %v", err.Error())
	}

	for _, a := range output.Accounts {
		if a.Currency == "USD" {
			assert.Equal(t, "99999", a.Balance)
			assert.Equal(t, "999", a.Hold)
			assert.Equal(t, "99000", a.Available)
		}
		if a.Currency == "ETH" {
			assert.Equal(t, "100010", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "100010", a.Available)
		}
	}

	_, err = ledgerClient.PartialReleaseHold(ctx, &ledger.PartialReleaseHoldRequest{
		OrderId:        orderId,
		RequestId:      "7795D23A-D273-442B-B588-CD718E11E2E5",
		SenderAmount:   "999",
		ReceiverAmount: "10",
	})

	if err != nil {
		t.Fatalf("unable to partial release hold")
	}

	output, err = ledgerClient.GetAccounts(ctx, &ledger.GetAccountsRequest{
		UserId: "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
	})

	if err != nil {
		t.Fatalf("complete transaction test: unable get accounts: %v", err.Error())
	}

	for _, a := range output.Accounts {
		if a.Currency == "USD" {
			assert.Equal(t, "99000", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "99000", a.Available)
		}
		if a.Currency == "ETH" {
			assert.Equal(t, "100020", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "100020", a.Available)
		}
	}

	_, err = ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         orderId,
		RequestId:       "E2586CC8-2879-48B2-83E2-F1F5CF87E248",
		FinalizedStatus: ledger.TransactionStatus_COMPLETE,
	})

	if err != nil {
		t.Fatalf("unable to complete transaction")
	}

	output, err = ledgerClient.GetAccounts(ctx, &ledger.GetAccountsRequest{
		UserId: "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
	})

	if err != nil {
		t.Fatalf("unable to get accounts")
	}

	for _, a := range output.Accounts {
		if a.Currency == "USD" {
			assert.Equal(t, "99000", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "99000", a.Available)
		}
		if a.Currency == "ETH" {
			assert.Equal(t, "100020", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "100020", a.Available)
		}
	}
}
