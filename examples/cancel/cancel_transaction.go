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
	"log"

	ledger "LedgerApp/protos/ledger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	ledgerClient := NewLedgerServiceClient("localhost:8443")

	transactionRes, err := ledgerClient.CreateTransaction(ctx, &ledger.CreateTransactionRequest{
		OrderId:         "456B4BF7-D975-4AED-B0F0-33FC1666F69B",
		SenderId:        "B183F5E2-B72A-4AA5-B7AE-95E0D548D84D",
		ReceiverId:      "9AA945E8-05FB-4D8E-88C9-1986F0813292",
		SenderAmount:    "10",
		TransactionType: ledger.TransactionType_TRANSFER,
		RequestId:       "33D81CC5-42AE-4BCE-86F5-915338A9FDFF",
	})

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to create transaction")
	}

	_, err = ledgerClient.PartialReleaseHold(ctx, &ledger.PartialReleaseHoldRequest{
		OrderId:        transactionRes.Transaction.Id,
		RequestId:      "11A40B5B-74AA-4CBD-A04B-AADC4F0487E8",
		SenderAmount:   "1",
		ReceiverAmount: "10",
	})

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to partial release hold")
	}

	_, err = ledgerClient.FinalizeTransaction(ctx, &ledger.FinalizeTransactionRequest{
		OrderId:         transactionRes.Transaction.Id,
		RequestId:       "E2586CC8-2879-48B2-83E2-F1F5CF87E248",
		FinalizedStatus: ledger.TransactionStatus_CANCELED,
	})

	if err != nil {
		log.Print(err)
		log.Fatalf("unable to complete transaction")
	}
}

func NewLedgerServiceClient(uri string) ledger.LedgerClient {
	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("unable to create ledger client")
	}

	return ledger.NewLedgerClient(conn)
}
