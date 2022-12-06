package test

import (
	"context"
	"testing"

	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
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
