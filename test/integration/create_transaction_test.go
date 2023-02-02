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

	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCreateTransactionAccountNotFound(t *testing.T) {
	ctx := context.Background()
	ledgerClient := newLedgerServiceClient(ctx, t)

	if _, err := ledgerClient.CreateTransaction(ctx, &ledger.CreateTransactionRequest{
		OrderId: "5AFD5F86-AF2C-45D8-8D92-105EC153A0C6",
		Sender: &ledger.Account{
			UserId:   "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
			Currency: "CELO",
		},
		Receiver: &ledger.Account{
			UserId:   "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
			Currency: "ETH",
		},
		TotalAmount:     "1000",
		FeeAmount:       &wrapperspb.StringValue{Value: "5"},
		TransactionType: ledger.TransactionType_TRANSACTION_TYPE_TRANSFER,
		RequestId:       &wrapperspb.StringValue{Value: "27AA0E33-FC64-4D80-A811-C9BB3416692A"},
	}); err == nil {
		t.Fatal("expected CreateTransaction with CELO account to fail")
	}
}
