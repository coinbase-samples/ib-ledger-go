/**
 * Copyright 2023-present Coinbase Global, Inc.
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

package qldb

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/amzn/ion-go/ion"
	"github.com/awslabs/amazon-qldb-driver-go/v3/qldbdriver"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
)

func CreateTransactionAndPlaceHold(
	ctx context.Context,
	t *model.QldbTransaction,
	amount *big.Int,
) error {
	_, err := Repo.Driver.Execute(
		ctx,
		func(txn qldbdriver.Transaction) (interface{}, error) {
			// Check to make sure Transaction wasn't already written
			result, err := txn.Execute(
				"SELECT * FROM Ledger WHERE id = ?",
				t.Id,
			)
			if err != nil {
				return nil, err
			}
			if result.Next(txn) {
				return nil, nil
			}

			sender, err := getAccountQldbTransaction(txn, t.Sender.Id)
			if err != nil {
				return nil, err
			}
			available, err := utils.IonDecimalToBigInt(sender.Available)
			if err != nil {
				return nil, err
			}

			if available.Cmp(amount) == -1 {
				return nil, &InsufficientBalanceError{}
			}

			// Write transaction
			_, err = txn.Execute("INSERT INTO Ledger ?", t)
			if err != nil {
				return nil, err
			}

			// Update Account Balance
			holdBalanceUpdate(txn, sender, amount, false)

			return nil, err
		})

	return err
}

func FinalizeTransactionAndReleaseHold(
	ctx context.Context,
	venueOrderId, status string) error {
	_, err := Repo.Driver.Execute(
		ctx,
		func(txn qldbdriver.Transaction) (interface{}, error) {
			// Retrieve transaction
			t, err := getTransactionQldbOperation(txn, venueOrderId)
			if err != nil {
				return nil, err
			}

			sender, err := getAccountQldbTransaction(txn, t.Sender.Id)
			if err != nil {
				return nil, err
			}

			h := &model.QldbHold{
				AccountId:  t.Hold.AccountId,
				Amount:     ion.NewDecimalInt(0),
				ReleasedAt: time.Now(),
				Released:   true,
				HoldUUID:   t.Hold.HoldUUID,
			}
			_, err = txn.Execute(
				"UPDATE Ledger AS t SET t.hold = ?, t.status = ? WHERE t.id = ?",
				h,
				status,
				t.Id,
			)

			if err != nil {
				return nil, err
			}

			txnHoldAmount, err := utils.IonDecimalToBigInt(t.Hold.Amount)
			if err != nil {
				return nil, err
			}

			err = holdBalanceUpdate(txn, sender, txnHoldAmount, true)

			return nil, err
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to finalize transaction and release hold - orderId: %s - status: %s - %w",
			venueOrderId,
			status,
			err,
		)
	}
	return nil
}

func GetTransaction(
	ctx context.Context,
	venueOrderId string,
) (*model.QldbTransaction, error) {
	data, err := Repo.Driver.Execute(ctx,
		func(txn qldbdriver.Transaction) (interface{}, error) {
			return getTransactionQldbOperation(txn, venueOrderId)
		},
	)
	if err != nil {
		return nil, err
	}
	if transaction, ok := data.(*model.QldbTransaction); ok {
		return transaction, nil
	} else {
		return nil, fmt.Errorf(
			"unable to cast data to transaction type: %v",
			data,
		)
	}
}

func getTransactionQldbOperation(
	txn qldbdriver.Transaction,
	venueOrderId string,
) (*model.QldbTransaction, error) {
	txnId := model.GenerateTransactionId(venueOrderId)
	result, err := txn.Execute("SELECT * FROM Ledger WHERE id = ?", txnId)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to execute get transaction - transactionId: %s - error: %w",
			venueOrderId,
			err,
		)
	}
	if !result.Next(txn) {
		return nil, fmt.Errorf(
			"failed to get result from get transaction - transactionId: %s - err: %w",
			venueOrderId,
			err,
		)
	}

	t := new(model.QldbTransaction)
	if err := ion.Unmarshal(result.GetCurrentData(), &t); err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal transaction result - transactionId: %s - err: %w",
			venueOrderId,
			err,
		)
	}
	return t, nil
}
