package test

import (
	"context"
	"testing"

	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestSuccessfulCancelTransaction(t *testing.T) {
	ctx := context.Background()

	ledgerClient := newLedgerServiceClient(ctx, t)

	orderId := "356798E8-2CE3-4D49-BD66-2BEC6C95AF4A"
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
		RequestId:       &wrapperspb.StringValue{Value: "E9115CD9-8E15-44B9-A3C1-4F5EC87FE12F"},
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

	finalizeResult, err := ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         orderId,
		RequestId:       "A23CE1FB-4221-4CB0-8E8F-8D5B7C3C4FA3",
		FinalizedStatus: ledger.TransactionStatus_CANCELED,
	})

	if err != nil {
		t.Fatalf("unable to fail transaction: %v", err.Error())
	}

	assert.True(t, finalizeResult.Successful)

	senderExpectedBalance = &ledger.AccountAndBalance{
		Currency:  "USD",
		Balance:   "100000",
		Hold:      "0",
		Available: "100000",
	}
	receiverExpectedBalance = &ledger.AccountAndBalance{
		Currency:  "ETH",
		Balance:   "100000",
		Hold:      "0",
		Available: "100000",
	}
	getTransactionBalancesAndConfirmTheyAreAsExpected(ledgerClient, ctx, t, senderExpectedBalance, receiverExpectedBalance)
}
