package test

import (
	"context"
	"strings"
	"testing"

	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func createTransactionAndConfirmHolds(l ledger.LedgerClient, ctx context.Context, t *testing.T, orderId string, requestId string) {
	transactionResult, err := l.CreateTransaction(ctx, &ledger.CreateTransactionRequest{
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
		RequestId:       &wrapperspb.StringValue{Value: requestId},
	})

	if err != nil {
		t.Fatalf("unable to create transaction: %v", err.Error())
	}

	assert.Equal(t, orderId, strings.ToUpper(transactionResult.Transaction.Id))
	assert.Equal(t, requestId, strings.ToUpper(transactionResult.Transaction.RequestId))

}

func partialReleaseHoldAndConfirmSuccessful(l ledger.LedgerClient, ctx context.Context, t *testing.T, pr *ledger.PartialReleaseHoldRequest) {
	tr, err := l.PartialReleaseHold(ctx, pr)

	if err != nil {
		t.Fatalf("unable to partial release hold")
	}

	assert.True(t, tr.Successful)
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
