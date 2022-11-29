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

package main

import (
	"context"
	"fmt"
	"log"

	ledger "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	ctx := context.Background()

	md := metadata.New(map[string]string{"x-route-id": "ledger"})

	conn, err := grpc.DialContext(
		metadata.NewOutgoingContext(context.TODO(), md),
		"localhost:8445",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to generate ledger connection")
	}

	ledgerClient := ledger.NewLedgerClient(conn)

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
		TotalAmount:     "1000",
		FeeAmount:       &wrapperspb.StringValue{Value: "5"},
		TransactionType: ledger.TransactionType_TRANSFER,
		RequestId:       &wrapperspb.StringValue{Value: "33D81CC5-42AE-4BCE-86F5-915338A9FDFF"},
	})

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to create transaction")
	}

	output, err := ledgerClient.GetAccounts(ctx, &ledger.GetAccountsRequest{
		UserId: "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
	})

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to get accounts")
	}

	for _, a := range output.Accounts {
		fmt.Println(fmt.Sprintf("accountId: %s - currency: %s - balance: %s, hold: %s, available: %s", a.AccountId, a.Currency, a.Balance, a.Hold, a.Available))
	}

	_, err = ledgerClient.PartialReleaseHold(ctx, &ledger.PartialReleaseHoldRequest{
		OrderId:         transactionRes.Transaction.Id,
		RequestId:       "11A40B5B-74AA-4CBD-A04B-AADC4F0487E8",
		SenderAmount:    "1",
		ReceiverAmount:  "10",
		VenueFeeAmount:  &wrapperspb.StringValue{Value: "1"},
		RetailFeeAmount: &wrapperspb.StringValue{Value: "1"},
	})

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to partial release hold")
	}

	_, err = ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         transactionRes.Transaction.Id,
		RequestId:       "E2586CC8-2879-48B2-83E2-F1F5CF87E248",
		FinalizedStatus: ledger.TransactionStatus_COMPLETE,
		SenderAmount:    &wrapperspb.StringValue{Value: "999"},
		ReceiverAmount:  &wrapperspb.StringValue{Value: "56"},
		VenueFeeAmount:  &wrapperspb.StringValue{Value: "1"},
		RetailFeeAmount: &wrapperspb.StringValue{Value: "2"},
	})

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to complete transaction")
	}

	output, err = ledgerClient.GetAccounts(ctx, &ledger.GetAccountsRequest{
		UserId: "620E62FD-DAF1-4738-84CE-1DBC4393ED29",
	})

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to get accounts")
	}

	for _, a := range output.Accounts {
		fmt.Println(fmt.Sprintf("accountId: %s - currency: %s - balance: %s, hold: %s, available: %s", a.AccountId, a.Currency, a.Balance, a.Hold, a.Available))
	}
}

func NewLedgerServiceClient(uri string) ledger.LedgerClient {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("unable to create ledger client")
	}

	return ledger.NewLedgerClient(conn)
}
