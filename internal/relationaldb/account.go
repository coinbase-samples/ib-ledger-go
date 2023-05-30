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

package relationaldb

import (
	"context"
	"fmt"

	"github.com/coinbase-samples/ib-ledger-go/internal/utils"

	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/google/uuid"
)

const (
	insertAccountSql = `
        INSERT INTO account (id, qldb_id, user_id, currency) 
        VALUES ($1, $2, $3, $4) 
        ON CONFLICT (id)
        DO NOTHING`

	selectAccountById = `SELECT id FROM account WHERE id = $1`

	insertAccountBalanceSql = `
        INSERT INTO account_balance (account_id, balance, hold, available, created_at, idem)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT ON CONSTRAINT balance_change
        DO NOTHING`

	getAllAccountsAndBalancesSql = `
        SELECT acc.id, acc.currency, ab.balance, ab.hold, ab.available, ab.created_at
        FROM (select id, currency FROM account WHERE user_id = $1) acc
        INNER JOIN
        (SELECT account_id, MAX(count) as max
        FROM account_balance
        WHERE account_id IN (select id
        FROM account
        WHERE user_id = $1)
        GROUP BY account_id) recent_balance
        ON acc.id = recent_balance.account_id
        INNER JOIN
        account_balance ab
        ON recent_balance.account_id = ab.account_id and recent_balance.max = ab.count
        GROUP BY ab.count, acc.id, acc.currency, ab.balance, ab.hold, ab.available, ab.created_at
        HAVING recent_balance.count = acc.count;`

	getAccountByIdSql = `SELECT * FROM account WHERE id = $1`
)

func insertAccount(ctx context.Context, account *model.Account) error {
	if err := Repo.Insert(
		ctx,
		insertAccountSql,
		account.Id.String(),
		account.QldbId,
		account.UserId.String(),
		account.Currency,
	); err != nil {
		return fmt.Errorf(
			"failed to insert account - id: %s - %w",
			account.Id,
			err,
		)
	}
	return nil
}

func InsertAccountBalance(
	ctx context.Context,
	account *model.QldbAccount,
) error {
	id := account.AccountUUID
	userId := account.UserId
	currency := account.Currency
	var accountIds []*uuid.UUID
	if err := Repo.Query(
		ctx,
		&accountIds,
		selectAccountById,
		id,
	); err != nil {
		return fmt.Errorf(
			"failed to retrieve account - accountId: %s - userId: %s - currency: %s - %w",
			id,
			userId,
			currency,
			err,
		)
	}

	// If first attempt doesn't find an accountId, create the new account
	if len(accountIds) == 0 {
		if err := Repo.Insert(
			ctx,
			insertAccountSql,
			id,
			account.Id,
			userId,
			currency,
		); err != nil {
			return fmt.Errorf(
				"failed to create account - accountId: %s - userId: %s - currency: %s - %w",
				id,
				userId,
				currency,
				err,
			)
		}
	}

	balance, err := utils.IonDecimalToBigInt(account.Balance)
	if err != nil {
		return fmt.Errorf(
			"failed account balance update - bad balance: %s - %w",
			account.Balance.String(),
			err,
		)
	}
	hold, _ := utils.IonDecimalToBigInt(account.Hold)
	if err != nil {
		return fmt.Errorf(
			"failed account balance update - bad hold: %s - %w",
			account.Hold.String(),
			err,
		)
	}
	available, _ := utils.IonDecimalToBigInt(account.Available)
	if err != nil {
		return fmt.Errorf(
			"failed account balance update - bad available: %s - %w",
			account.Available.String(),
			err,
		)
	}
	idem := utils.GenerateIdemString(
		id,
		balance.String(),
		hold.String(),
		available.String(),
		account.UpdatedAt.String(),
	)

	if err := Repo.Insert(
		ctx,
		insertAccountBalanceSql,
		id,
		balance.String(),
		hold.String(),
		available.String(),
		account.UpdatedAt,
		idem,
	); err != nil {
		return fmt.Errorf(
			"unable to insert account balance - accountId: %s - %w",
			id,
			err,
		)
	}
	return nil
}

func GetAllAccountsAndMostRecentBalances(
	ctx context.Context,
	userId string,
) ([]*model.AccountBalance, error) {
	var data []*model.AccountBalance

	if err := Repo.Query(
		ctx,
		&data,
		getAllAccountsAndBalancesSql,
		userId,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to get account balances - userId: %s - %w",
			userId,
			err,
		)
	}
	return data, nil
}

func GetAccount(
	ctx context.Context,
	accountId string,
) (*model.Account, error) {
	var data *model.Account
	if err := Repo.Query(
		ctx,
		&data,
		getAccountByIdSql,
		accountId,
	); err != nil {
		return nil, fmt.Errorf(
			"unable to get account - id: %s - %w",
			accountId,
			err,
		)
	}
	return data, nil
}
