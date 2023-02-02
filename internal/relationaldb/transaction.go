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
	"strings"

	"github.com/coinbase-samples/ib-ledger-go/internal/model"
)

const (
	insertTransactionSql = `
        INSERT INTO transaction (id, qldb_id, sender_id, receiver_id, created_at, finalized_at, transaction_status, transaction_type)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (id)
        DO UPDATE SET 
            finalized_at = $6, transaction_status = $7
        WHERE EXCLUDED.finalized_at IS NOT NULL`

	insertEntriesSql = `
        INSERT INTO entry (id, account_id, venue_order_id, fill_id, amount, direction, created_at)
        VALUES($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14)
        ON CONFLICT (id)
        DO NOTHING`

	insertHoldSql = `
        INSERT INTO hold (id, account_id, venue_order_id, amount, created_at, released_at, released)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id)
        DO UPDATE SET 
            amount = EXCLUDED.amount, released_at = EXCLUDED.released_at, released = EXCLUDED.released
        WHERE hold.amount > EXCLUDED.amount OR hold.released != EXCLUDED.released`

	deleteEntrySql = `DELETE FROM entry WHERE id = $1`
)

func InsertTransaction(
	ctx context.Context,
	t *model.Transaction,
	sender, receiver *model.Account,
) error {
	if err := checkIfAccountExistsCreateIfNotFound(
		ctx,
		sender,
	); err != nil {
		return fmt.Errorf(
			"failed to get sender account - id: %s - %w",
			sender.Id.String(),
			err,
		)
	}

	if err := checkIfAccountExistsCreateIfNotFound(
		ctx,
		receiver,
	); err != nil {
		return fmt.Errorf(
			"failed to get receiver account - id: %s - %w",
			receiver.Id.String(),
			err,
		)
	}

	if err := Repo.Insert(
		ctx,
		insertTransactionSql,
		t.Id.String(),
		t.QldbId,
		t.Sender.String(),
		t.Receiver.String(),
		t.CreatedAt,
		t.FinalizedAt,
		strings.ToUpper(t.TransactionStatus),
		t.TransactionType,
	); err != nil {
		return fmt.Errorf(
			"failed to insert transaction - id: %s - %w",
			t.Id.String(),
			err,
		)
	}
	return nil
}

func InsertEntry(ctx context.Context, a, b *model.Entry) error {
	appendArray := append([]string(nil), a.ToRow()...)
	appendArray = append(appendArray, b.ToRow()...)
	if err := Repo.Insert(
		ctx,
		insertEntriesSql,
		appendArray,
	); err != nil {
		return fmt.Errorf(
			"failed to insert entries - entry1: %s - entry2: %s - %w",
			a.Id.String(),
			b.Id.String(),
			err,
		)
	}
	return nil
}

func DeleteEntry(ctx context.Context, id string) error {
	if err := Repo.Insert(ctx, deleteEntrySql, id); err != nil {
		return fmt.Errorf(
			"unable to delete entry - id: %s - %w",
			id,
			err,
		)
	}
	return nil
}

func UpsertHold(ctx context.Context, hold *model.Hold) error {
	if err := Repo.Insert(
		ctx,
		insertHoldSql,
		hold.Id,
		hold.AccountId,
		hold.VenueOrderId,
		hold.Amount,
		hold.CreatedAt,
		hold.ReleasedAt,
		hold.Released,
	); err != nil {
		return fmt.Errorf("unable to upsert hold - holdId: %s - transactionId: %s - %w",
			hold.Id.String(),
			hold.VenueOrderId.String(),
			err,
		)
	}
	return nil
}

func checkIfAccountExistsCreateIfNotFound(
	ctx context.Context,
	a *model.Account,
) error {
	var retrievedId []*string
	if err := Repo.Query(ctx,
		&retrievedId,
		selectAccountById,
		a.Id.String(),
	); err != nil {
		return fmt.Errorf(
			"failed to retrieve account - id: %s - %w",
			a.Id.String(),
			err,
		)
	}

	// If first attempt doesn't find an accountId, create the new account
	if len(retrievedId) == 0 {
		if err := insertAccount(ctx, a); err != nil {
			return fmt.Errorf(
				"unable to insert account - id: %s - %w",
				a.Id.String(),
				err,
			)
		}
	}
	return nil
}
