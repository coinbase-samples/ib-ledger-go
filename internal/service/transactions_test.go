package service

import (
	"context"
	"fmt"
	"testing"

	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBadFinalizeTransactionType(t *testing.T) {
    orderId := "5E45DE98-7462-45C7-8801-37474B3F5EE4"
    service := NewTestService()
    request := &api.FinalizeTransactionRequest{
        OrderId: orderId,
        SenderAmount: &wrapperspb.StringValue{Value: "1"},
        FinalizedStatus: api.TransactionStatus_PENDING,
    }

    _, err := service.FinalizeTransaction(context.TODO(), request)

    if err == nil {
        t.Fatal("expected error returned from function invocation")
    }

    if err.Error() != fmt.Sprintf("finalize transaction: unable to finalize pending transaction - transaction: %v", orderId) {
        t.Fatalf("unexpected error message, received: %v", err.Error())
    }
} 

func TestTransactionNotFoundCompleteTransaction(t *testing.T) {
    service := NewTestService()
    request := &api.FinalizeTransactionRequest{
        OrderId: CompleteTransactionUuidNotFound,
        SenderAmount: &wrapperspb.StringValue{ Value: "1" },
        FinalizedStatus: api.TransactionStatus_COMPLETE,
    }
    
    _, err := service.FinalizeTransaction(context.TODO(), request)

    if err == nil {
        t.Fatal("expected transaction not found error returned from complete transaction invocation")
    }

    if err.Error() != "transaction not found" {
        t.Fatalf("unexpected error message from complete transaction, received: %v", err.Error())
    }
}

func TestSuccessfulCompleteTransaction(t *testing.T) {
    orderId := "5E45DE98-7462-45C7-8801-37474B3F5EE4"
    service := NewTestService()
    request := &api.FinalizeTransactionRequest{
        OrderId: orderId,
        SenderAmount: &wrapperspb.StringValue{ Value: "1" },
        FinalizedStatus: api.TransactionStatus_COMPLETE,
    }
    
    _, err := service.FinalizeTransaction(context.TODO(), request)

    if err != nil {
        t.Fatalf("expected successful execution of complete transaction, received err: %v", err.Error())
    }
}

func TestTransactionNotFoundFailTransaction(t *testing.T) {
    service := NewTestService()
    request := &api.FinalizeTransactionRequest{
        OrderId: FailTransactionUuidNotFound,
        SenderAmount: &wrapperspb.StringValue{ Value: "1" },
        FinalizedStatus: api.TransactionStatus_FAILED,
    }
    
    _, err := service.FinalizeTransaction(context.TODO(), request)

    if err == nil {
        t.Fatal("expected transaction not found error returned from fail transaction invocation")
    }

    if err.Error() != "transaction not found" {
        t.Fatalf("unexpected error message from fail transaction, received: %v", err.Error())
    }
}

func TestSuccessfulFailTransaction(t *testing.T) {
    orderId := "5E45DE98-7462-45C7-8801-37474B3F5EE4"
    service := NewTestService()
    request := &api.FinalizeTransactionRequest{
        OrderId: orderId,
        SenderAmount: &wrapperspb.StringValue{ Value: "1" },
        FinalizedStatus: api.TransactionStatus_FAILED,
    }
    
    _, err := service.FinalizeTransaction(context.TODO(), request)

    if err != nil {
        t.Fatalf("expected successful execution of fail transaction, received err: %v", err.Error())
    }
}

func TestTransactionNotFoundCancelTransaction(t *testing.T) {
    service := NewTestService()
    request := &api.FinalizeTransactionRequest{
        OrderId: CancelTransactionUuidNotFound,
        SenderAmount: &wrapperspb.StringValue{ Value: "1" },
        FinalizedStatus: api.TransactionStatus_CANCELED,
    }
    
    _, err := service.FinalizeTransaction(context.TODO(), request)

    if err == nil {
        t.Fatal("expected transaction not found error returned from cancel transaction invocation")
    }

    if err.Error() != "transaction not found" {
        t.Fatalf("unexpected error message from cancel transaction, received: %v", err.Error())
    }
}

func TestSuccessfulCancelTransaction(t *testing.T) {
    orderId := "5E45DE98-7462-45C7-8801-37474B3F5EE4"
    service := NewTestService()
    request := &api.FinalizeTransactionRequest{
        OrderId: orderId,
        SenderAmount: &wrapperspb.StringValue{ Value: "1" },
        FinalizedStatus: api.TransactionStatus_CANCELED,
    }
    
    _, err := service.FinalizeTransaction(context.TODO(), request)

    if err != nil {
        t.Fatalf("expected successful execution of cancel transaction, received err: %v", err.Error())
    }
}
