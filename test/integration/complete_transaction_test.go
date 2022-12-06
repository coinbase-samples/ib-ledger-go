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

func TestCompleteTransaction(t *testing.T) {
	ctx := context.Background()

	ledgerClient := newLedgerServiceClient(ctx, t)

	const totalAmount = "1000"
	const feeAmount = "5"

	transactionRes, err := ledgerClient.CreateTransaction(ctx, &ledger.CreateTransactionRequest{
		OrderId: "456B4BF7-D975-4AED-B0F0-33FC1666F69B",
		Sender: &ledger.Account{
			UserId:   "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
			Currency: "USD",
		},
		Receiver: &ledger.Account{
			UserId:   "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
			Currency: "ETH",
		},
		TotalAmount:     totalAmount,
		FeeAmount:       &wrapperspb.StringValue{Value: feeAmount},
		TransactionType: ledger.TransactionType_TRANSFER,
		RequestId:       &wrapperspb.StringValue{Value: "33D81CC5-42AE-4BCE-86F5-915338A9FDFF"},
	})

	if err != nil {
		t.Fatalf("complete transaction test: unable to create transaction: %v", err.Error())
	}

	output, err := ledgerClient.GetAccounts(ctx, &ledger.GetAccountsRequest{
		UserId: "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
	})

	if err != nil {
		t.Fatalf("complete transaction test: unable get accounts: %v", err.Error())
	}

	isUsdAccountPresent := false
	isEthAccountPresent := false
	for _, a := range output.Accounts {
		if a.Currency == "USD" {
			isUsdAccountPresent = true
			assert.Equal(t, "100000", a.Balance)
			assert.Equal(t, "1005", a.Hold)
			assert.Equal(t, "98995", a.Available)
		}
		if a.Currency == "ETH" {
			isEthAccountPresent = true
			assert.Equal(t, "100000", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "100000", a.Available)
		}
	}
	assert.True(t, isUsdAccountPresent)
	assert.True(t, isEthAccountPresent)

	_, err = ledgerClient.PartialReleaseHold(ctx, &ledger.PartialReleaseHoldRequest{
		OrderId:         transactionRes.Transaction.Id,
		RequestId:       "11A40B5B-74AA-4CBD-A04B-AADC4F0487E8",
		SenderAmount:    "1",
		ReceiverAmount:  "10",
		VenueFeeAmount:  &wrapperspb.StringValue{Value: "1"},
		RetailFeeAmount: &wrapperspb.StringValue{Value: "1"},
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
			assert.Equal(t, "99997", a.Balance)
			assert.Equal(t, "1002", a.Hold)
			assert.Equal(t, "98995", a.Available)
		}
		if a.Currency == "ETH" {
			assert.Equal(t, "100010", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "100010", a.Available)
		}
	}

	_, err = ledgerClient.PartialReleaseHold(ctx, &ledger.PartialReleaseHoldRequest{
		OrderId:         transactionRes.Transaction.Id,
		RequestId:       "7795D23A-D273-442B-B588-CD718E11E2E5",
		SenderAmount:    "999",
		ReceiverAmount:  "10",
		VenueFeeAmount:  &wrapperspb.StringValue{Value: "2"},
		RetailFeeAmount: &wrapperspb.StringValue{Value: "1"},
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
			assert.Equal(t, "98995", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "98995", a.Available)
		}
		if a.Currency == "ETH" {
			assert.Equal(t, "100020", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "100020", a.Available)
		}
	}

	_, err = ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         transactionRes.Transaction.Id,
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
			assert.Equal(t, "98995", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "98995", a.Available)
		}
		if a.Currency == "ETH" {
			assert.Equal(t, "100020", a.Balance)
			assert.Equal(t, "0", a.Hold)
			assert.Equal(t, "100020", a.Available)
		}
	}
}

