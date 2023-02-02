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
	"time"

	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func newLedgerServiceClient(ctx context.Context, t *testing.T) ledger.LedgerClient {
	md := metadata.New(map[string]string{"x-route-id": "ledger"})

	conn, err := grpc.DialContext(
		metadata.NewOutgoingContext(ctx, md),
		"localhost:8445",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		t.Fatalf("complete transaction test: unable to generate ledger connection: %v", err.Error())
	}

	return ledger.NewLedgerClient(conn)
}

func setupAccounts(t *testing.T, ctx context.Context, l ledger.LedgerClient, userId string) {
	if _, err := l.InitializeAccount(ctx, &ledger.InitializeAccountRequest{
		UserId:         userId,
		Currency:       "USD",
		InitialBalance: "100000",
	}); err != nil {
		t.Fatalf("unable to initialize USD account: %v", err)
	}
	if _, err := l.InitializeAccount(ctx, &ledger.InitializeAccountRequest{
		UserId:         userId,
		Currency:       "ETH",
		InitialBalance: "0",
	}); err != nil {
		t.Fatalf("unable to initialize ETH account: %v", err)
	}
	time.Sleep(5000)
}

func getAccountsAndConfirmBalances(t *testing.T, ctx context.Context, l ledger.LedgerClient, userId string, account1Expected *model.QldbAccount, account2Expected *model.QldbAccount) {
	res, err := l.GetAccounts(ctx, &ledger.GetAccountsRequest{UserId: userId})
	assert.Nil(t, err)

	account1equal := false
	account2equal := false
	for _, v := range res.Accounts {
		if v.AccountId == account1Expected.Id {
			account1equal = account1Expected.Equal(v)
		} else if v.AccountId == account2Expected.Id {
			account2equal = account2Expected.Equal(v)
		}
	}
	assert.True(t, account1equal)
	assert.True(t, account2equal)
}

func createTransaction(l ledger.LedgerClient, ctx context.Context, t *testing.T, req *ledger.CreateTransactionRequest) {
	transactionResult, err := l.CreateTransaction(ctx, req)

	if err != nil {
		t.Fatalf("unable to create transaction: %v", err.Error())
	}

	assert.NotNil(t, transactionResult)
}

func postFill(l ledger.LedgerClient, ctx context.Context, t *testing.T, pr *ledger.PostFillRequest) {
	_, err := l.PostFill(ctx, pr)

	if err != nil {
		t.Fatalf("unable to partial release hold")
	}
}

func getTransactionBalancesAndConfirmTheyAreAsExpected(
	l ledger.LedgerClient,
	ctx context.Context,
	t *testing.T,
	senderBalance *ledger.AccountAndBalance,
	receiverBalance *ledger.AccountAndBalance) {
	output, err := l.GetAccounts(ctx, &ledger.GetAccountsRequest{
		UserId: "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
	})

	if err != nil {
		t.Fatalf("unable get accounts: %v", err.Error())
	}

	isSenderAccountBalancePresent := false
	isReceiverAccountBalancePresent := false
	for _, a := range output.Accounts {
		if a.Currency == senderBalance.Currency {
			isSenderAccountBalancePresent = true
			assert.True(t, accountBalancesAreEqual(senderBalance, a))
		}
		if a.Currency == receiverBalance.Currency {
			isReceiverAccountBalancePresent = true
			assert.True(t, accountBalancesAreEqual(receiverBalance, a))
		}
	}
	assert.True(t, isSenderAccountBalancePresent)
	assert.True(t, isReceiverAccountBalancePresent)
}

func accountBalancesAreEqual(a *ledger.AccountAndBalance, b *ledger.AccountAndBalance) bool {
	return a.Balance == b.Balance && a.Hold == b.Hold && a.Available == b.Available
}
