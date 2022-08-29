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
