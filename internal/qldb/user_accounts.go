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

	"github.com/amzn/ion-go/ion"
	"github.com/awslabs/amazon-qldb-driver-go/v3/qldbdriver"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
)

const (
	getAccountsByUserIdSql = `SELECT * FROM Accounts WHERE userId = ?`
	getAccountByIdSql      = `SELECT * FROM Accounts WHERE id = ?`
)

func GetUserAccounts(
	ctx context.Context,
	userId string,
) ([]*model.QldbAccount, error) {
	data, err := Repo.Driver.Execute(
		ctx,
		func(txn qldbdriver.Transaction) (interface{}, error) {
			result, err := txn.Execute(getAccountsByUserIdSql, userId)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to execute get accounts transaction - %w",
					err,
				)
			}
			var output []*model.QldbAccount
			for result.Next(txn) {
				var temp *model.QldbAccount
				if err := ion.Unmarshal(
					result.GetCurrentData(), &temp); err != nil {
					return nil, fmt.Errorf(
						"failed to unmarshal account data - %w",
						err,
					)
				}
				output = append(output, temp)
			}
			return output, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get accounts by userId from qldb - userId: %s - %w",
			userId,
			err,
		)
	}
	account, ok := data.([]*model.QldbAccount)
	if !ok {
		return nil, fmt.Errorf(
			"unable to cast output data as account data array - userId: %s",
			userId,
		)
	}
	return account, nil
}

func GetAccount(
	ctx context.Context,
	accountId string,
) (*model.QldbAccount, error) {
	var res interface{}
	var err error

	res, err = Repo.Driver.Execute(
		ctx,
		func(txn qldbdriver.Transaction) (interface{}, error) {
			result, err := txn.Execute(getAccountByIdSql, accountId)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to execute get account sql - %w",
					err,
				)
			}

			if !result.Next(txn) {
				return nil, &AccountNotFoundError{Id: accountId}
			}

			ionBinary := result.GetCurrentData()

			acct := new(model.QldbAccount)
			if err := ion.Unmarshal(ionBinary, acct); err != nil {
				return nil, fmt.Errorf(
					"unable to unmarshal to account: %w",
					err,
				)
			}

			return acct, err
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get account by id - accountId: %s - %w",
			accountId,
			err,
		)
	}

	account, ok := res.(*model.QldbAccount)
	if !ok {
		return nil, fmt.Errorf(
			"unable to cast output data as account - accountId: %s",
			accountId,
		)
	}
	return account, nil
}
