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
	"github.com/google/uuid"
)

const (
	selectAccountSql = "SELECT * FROM Accounts WHERE id = ?"
)

func CreateAccountAndInitializeBalance(
	ctx context.Context,
	userId, currency string,
	initialBalance *big.Int,
) error {
	_, err := Repo.Driver.Execute(ctx,
		func(txn qldbdriver.Transaction) (interface{}, error) {
			accountId := model.GenerateAccountId(userId, currency)
			result, err := txn.Execute(selectAccountSql, accountId)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to select account - id: %s - userId: %s - currency: %s - err: %w",
					accountId,
					userId,
					currency,
					err,
				)
			}
			if result.Next(txn) {
				return nil, nil
			}

			initialDecimal := ion.NewDecimal(initialBalance, 0, false)
			a := &model.QldbAccount{
				Id:          model.GenerateAccountId(userId, currency),
				Currency:    currency,
				UserId:      userId,
				Balance:     initialDecimal,
				Hold:        ion.MustParseDecimal("0"),
				Available:   initialDecimal,
				UpdatedAt:   time.Now(),
				AccountUUID: uuid.New().String(),
			}

			_, err = txn.Execute("INSERT INTO Accounts ?", a)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to insert account - userId: %s - currency: %s - err: %w",
					userId,
					currency,
					err,
				)
			}

			return nil, nil
		},
	)

	return err
}

func getAccountQldbTransaction(
	txn qldbdriver.Transaction,
	accountId string,
) (*model.QldbAccount, error) {
	var account *model.QldbAccount
	result, err := txn.Execute(selectAccountSql, accountId)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to execute get account transaction - accountId: %s - error: %w",
			accountId,
			err,
		)
	}
	if result.Next(txn) {
		if err := ion.Unmarshal(
			result.GetCurrentData(), &account); err != nil {
			return nil, fmt.Errorf(
				"failed to unmarshal account data - accountId: %s - error: %w",
				accountId,
				err,
			)
		}
	} else {
		return nil, &AccountNotFoundError{
			Id: accountId,
		}
	}
	if result.Err() != nil {
		return nil, fmt.Errorf(
			"failed to get account result - accountId: %s - error: %w",
			accountId,
			err,
		)
	}
	return account, nil
}

func creditAccountUpdate(
	txn qldbdriver.Transaction,
	account *model.QldbAccount,
	amount *big.Int,
) error {
	balance, err := utils.IonDecimalToBigInt(account.Balance)
	if err != nil {
		return err
	}

	balance.Add(balance, amount)

	hold, err := utils.IonDecimalToBigInt(account.Hold)
	if err != nil {
		return err
	}

	available, err := utils.IonDecimalToBigInt(account.Available)
	if err != nil {
		return err
	}
	available.Sub(balance, hold)

	_, err = txn.Execute(`
        UPDATE Accounts AS a 
        SET a.balance = ?, a.available = ?, a.updatedAt = ? 
        WHERE a.id = ?`,
		ion.NewDecimal(balance, 0, false),
		ion.NewDecimal(available, 0, false),
		time.Now(),
		account.Id,
	)
	return err
}

func debitAccountUpdate(
	txn qldbdriver.Transaction,
	account *model.QldbAccount,
	amount *big.Int,
	commission *big.Int,
	feeAmount *big.Int,
) error {

	hold, err := utils.IonDecimalToBigInt(account.Hold)
	if err != nil {
		return fmt.Errorf(
			"failed to update debit account balance - accountId: %s - error: %w",
			account.Id,
			err,
		)
	}

	hold.Sub(hold, amount)
	hold.Sub(hold, commission)
	hold.Sub(hold, feeAmount)

	balance, err := utils.IonDecimalToBigInt(account.Balance)
	if err != nil {
		return fmt.Errorf(
			"failed to update debit account balance - accountId: %s - error: %w",
			account.Id,
			err,
		)
	}
	balance.Sub(balance, amount)
	balance.Sub(balance, commission)
	balance.Sub(balance, feeAmount)

	var available *big.Int
	available.Sub(balance, hold)

	if _, err := txn.Execute(`
        UPDATE Accounts AS a 
        SET a.balance = ?, a.hold = ?, a.available = ?, a.updatedAt = ? 
        WHERE a.id = ?`,
		ion.NewDecimal(balance, 0, false),
		ion.NewDecimal(hold, 0, false),
		ion.NewDecimal(available, 0, false),
		time.Now(),
		account.Id,
	); err != nil {
		return fmt.Errorf(
			"failed debit balance update - accountId: %s - error: %w",
			account.Id,
			err,
		)
	}
	return nil
}

func holdBalanceUpdate(
	txn qldbdriver.Transaction,
	account *model.QldbAccount,
	amount *big.Int,
	isSub bool,
) error {

	hold, err := utils.IonDecimalToBigInt(account.Hold)
	if err != nil {
		return fmt.Errorf(
			"failed to update credit account balance - accountId: %s - error: %w",
			account.Id,
			err,
		)
	}

	if isSub {
		hold.Sub(hold, amount)
	} else {
		hold.Add(hold, amount)
	}

	balance, err := utils.IonDecimalToBigInt(account.Balance)
	if err != nil {
		return fmt.Errorf(
			"failed to update credit account balance - accountId: %s - error: %w",
			account.Id,
			err,
		)
	}

	available, err := utils.IonDecimalToBigInt(account.Available)
	if err != nil {
		return fmt.Errorf(
			"failed to update credit account balance - accountId: %s - error: %w",
			account.Id,
			err,
		)
	}

	available.Sub(balance, hold)

	if _, err := txn.Execute(`
        UPDATE Accounts AS a 
        SET a.hold = ?, a.available = ?, a.updatedAt = ? 
        WHERE a.id = ?`,
		ion.NewDecimal(hold, 0, false),
		ion.NewDecimal(available, 0, false),
		time.Now(),
		account.Id,
	); err != nil {
		return fmt.Errorf(
			"failed credit balance update - accountId: %s - error: %w",
			account.Id,
			err,
		)
	}
	return nil
}
