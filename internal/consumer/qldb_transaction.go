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

package consumer

import (
	"context"
	"fmt"

	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/coinbase-samples/ib-ledger-go/internal/relationaldb"
	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
	"github.com/google/uuid"
)

func writeQldbTransaction(ctx context.Context, t *model.QldbTransaction) error {

	transaction, err := t.ConvertToPostgresTransaction()
	if err != nil {
		return err
	}

	senderAccount, err := t.Sender.ConvertToPostgresAccount()
	if err != nil {
		return fmt.Errorf(
			"unable to convert QldbCoreAccount: %v to Postgres Account: %w",
			t.Sender,
			err,
		)
	}

	receiverAccount, err := t.Receiver.ConvertToPostgresAccount()
	if err != nil {
		return fmt.Errorf(
			"unable to convert QldbCoreAccount: %v to Postgres Account: %w",
			t.Receiver,
			err,
		)
	}
	if err := relationaldb.InsertTransaction(
		ctx,
		transaction,
		senderAccount,
		receiverAccount,
	); err != nil {
		return fmt.Errorf(
			"unable to insert transaction to postgres database: %w",
			err,
		)
	}

	holdAmount, err := utils.IonDecimalToBigInt(t.Hold.Amount)
	if err != nil {
		return err
	}

	holdId, err := uuid.Parse(t.Hold.HoldUUID)
	if err != nil {
		return fmt.Errorf("unable to parse hold id into valid uuid: %w", err)
	}
	h := &model.Hold{
		Id:           holdId,
		AccountId:    senderAccount.Id,
		Amount:       holdAmount.String(),
		VenueOrderId: transaction.Id,
		CreatedAt:    t.CreatedAt,
		Released:     t.Hold.Released,
		ReleasedAt:   t.Hold.ReleasedAt,
	}
	err = relationaldb.UpsertHold(ctx, h)
	return err
}
