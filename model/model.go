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

package model

import (
	"time"

	"github.com/google/uuid"
)

type InitializeAccountResult struct {
	Id          uuid.UUID `db:"id"`
	PortfolioId uuid.UUID `db:"portfolio_id"`
	UserId      uuid.UUID `db:"user_id"`
	Currency    string    `db:"currency"`
	CreatedAt   time.Time `db:"created_at"`
	Balance     int64     `db:"balance"`
	Hold        int64     `db:"hold"`
	Available   int64     `db:"available"`
}

type CreateTransactionResult struct {
	Id              uuid.UUID `db:"id"`
	SenderId        uuid.UUID `db:"sender_id"`
	ReceiverId      uuid.UUID `db:"receiver_id"`
	RequestId       uuid.UUID `db:"request_id"`
	TransactionType string    `db:"transaction_type"`
	CreatedAt       time.Time `db:"created_at"`
}

type TransactionResult struct {
	HoldId            uuid.UUID `db:"hold_id"`
	SenderEntryId     uuid.UUID `db:"sender_entry_id"`
	ReceiverEntryId   uuid.UUID `db:"receiver_entry_id"`
	SenderBalanceId   uuid.UUID `db:"sender_balance_id"`
	ReceiverBalanceId uuid.UUID `db:"receiver_balance_id"`
}

type GetAccountResult struct {
	AccountId uuid.UUID `db:"account_id"`
	Currency  string    `db:"currency"`
	Balance   string    `db:"balance"`
	Hold      string    `db:"hold"`
	Available string    `db:"available"`
	CreatedAt time.Time `db:"created_at"`
}
