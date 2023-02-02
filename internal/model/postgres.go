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

package model

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id       uuid.UUID `db:"id"`
	QldbId   string    `db:"qldb_id"`
	UserId   uuid.UUID `db:"user_id"`
	Currency string    `db:"currency"`
}

type AccountBalance struct {
	Id        uuid.UUID `db:"id"`
	AccountId uuid.UUID `db:"account_id"`
	Balance   string    `db:"balance"`
	Hold      string    `db:"hold"`
	Available string    `db:"available"`
	CreatedAt time.Time `db:"created_at"`
	RequestId uuid.UUID `db:"request_id"`
	Idem      string    `db:"idem"`
}

type Transaction struct {
	Id                uuid.UUID `db:"id"`
	QldbId            string    `db:"qldb_id"`
	Sender            uuid.UUID `db:"sender_id"`
	Receiver          uuid.UUID `db:"receiver_id"`
	CreatedAt         time.Time `db:"created_at"`
	FinalizedAt       time.Time `db:"finalized_at"`
	TransactionStatus string    `db:"transaction_status"`
	TransactionType   string    `db:"transaction_type"`
}

type Entry struct {
	Id           uuid.UUID `db:"id"`
	AccountId    uuid.UUID `db:"account_id"`
	VenueOrderId uuid.UUID `db:"venue_order_id"`
	FillId       uuid.UUID `db:"fill_id"`
	Amount       string    `db:"amount"`
	Direction    string    `db:"direction"`
	CreatedAt    time.Time `db:"created_at"`
}

func (e *Entry) ToRow() []string {
	return []string{
		e.Id.String(),
		e.AccountId.String(),
		e.VenueOrderId.String(),
		e.FillId.String(),
		e.Amount,
		e.Direction,
		e.CreatedAt.String(),
	}
}

type Hold struct {
	Id           uuid.UUID `db:"id"`
	AccountId    uuid.UUID `db:"account_id"`
	VenueOrderId uuid.UUID `db:"venue_order_id"`
	Amount       string    `db:"amount"`
	CreatedAt    time.Time `db:"created_at"`
	ReleasedAt   time.Time `db:"released_at"`
	Released     bool      `db:"released"`
	Idem         string    `db:"idem"`
}
