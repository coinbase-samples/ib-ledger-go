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
	"fmt"
	"time"

	"github.com/amzn/ion-go/ion"
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/google/uuid"
)

type QldbAccount struct {
	Id          string       `ion:"id"`
	Currency    string       `ion:"currency"`
	UserId      string       `ion:"userId"`
	Balance     *ion.Decimal `ion:"balance"`
	Hold        *ion.Decimal `ion:"hold"`
	Available   *ion.Decimal `ion:"available"`
	UpdatedAt   time.Time    `ion:"updatedAt"`
	AccountUUID string       `ion:"accountUUID"`
}

func (q *QldbAccount) GetCoreAccount() *QldbCoreAccount {
	return &QldbCoreAccount{
		Id:          q.Id,
		Currency:    q.Currency,
		UserId:      q.UserId,
		AccountUUID: q.AccountUUID,
	}
}

func (q *QldbAccount) ConvertToPostgresAccount() (*Account, error) {
	id, err := uuid.Parse(q.AccountUUID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse accountId into uuid: %w", err)
	}
	userId, err := uuid.Parse(q.UserId)
	if err != nil {
		return nil, fmt.Errorf("unable to parse userId into uuid: %w", err)
	}
	return &Account{
		Id:       id,
		QldbId:   q.Id,
		UserId:   userId,
		Currency: q.Currency,
	}, nil
}

func (q *QldbAccount) Equal(i *api.AccountAndBalance) bool {
	inputBalance, err := ion.ParseDecimal(i.Balance)
	if err != nil {
		return false
	}
	inputHold, err := ion.ParseDecimal(i.Hold)
	if err != nil {
		return false
	}
	inputAvailable, err := ion.ParseDecimal(i.Available)
	if err != nil {
		return false
	}
	return (q.Id == i.AccountId &&
		q.Currency == i.Currency &&
		q.Balance.Equal(inputBalance) &&
		q.Hold.Equal(inputHold) &&
		q.Available.Equal(inputAvailable))
}

type QldbTransaction struct {
	Id              string           `ion:"id"`
	VenueOrderId    string           `ion:"venueOrderId"`
	Sender          *QldbCoreAccount `ion:"sender"`
	Receiver        *QldbCoreAccount `ion:"receiver"`
	CreatedAt       time.Time        `ion:"createdAt"`
	UpdatedAt       time.Time        `ion:"updatedAt"`
	Status          string           `ion:"status"`
	TransactionType string           `ion:"transactionType"`
	Hold            *QldbHold        `ion:"hold"`
}

func (q *QldbTransaction) ConvertToPostgresTransaction() (*Transaction, error) {
	id, err := uuid.Parse(q.VenueOrderId)
	if err != nil {
		return nil, fmt.Errorf(
			"failed conversion to postgres transaction - bad orderId: %s - err: %w",
			q.VenueOrderId,
			err,
		)
	}
	senderId, err := uuid.Parse(q.Sender.AccountUUID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed conversion to postgres transaction - bad senderId: %s - err: %w",
			q.Sender.AccountUUID,
			err,
		)
	}
	receiverId, err := uuid.Parse(q.Receiver.AccountUUID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed conversion to postgres transaction - bad receiverId: %s - err: %w",
			q.Receiver.AccountUUID,
			err,
		)
	}
	return &Transaction{
		Id:                id,
		QldbId:            q.Id,
		Sender:            senderId,
		Receiver:          receiverId,
		CreatedAt:         q.CreatedAt,
		TransactionStatus: q.Status,
		TransactionType:   q.TransactionType,
		FinalizedAt:       q.UpdatedAt,
	}, nil
}

type QldbCoreAccount struct {
	Id          string `ion:"id"`
	Currency    string `ion:"currency"`
	UserId      string `ion:"userId"`
	AccountUUID string `ion:"accountUUID"`
}

func (q *QldbCoreAccount) ConvertToPostgresAccount() (*Account, error) {
	id, err := uuid.Parse(q.AccountUUID)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to convert to postgres account - bad accountId: %s - err: %w",
			q.AccountUUID,
			err,
		)
	}
	userId, err := uuid.Parse(q.UserId)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to convert to postgres account - bad userId: %s - err: %w",
			q.UserId,
			err,
		)
	}
	return &Account{
		Id:       id,
		QldbId:   q.Id,
		UserId:   userId,
		Currency: q.Currency,
	}, nil
}

type QldbFill struct {
	Id             string            `ion:"id"`
	ProductId      string            `ion:"productId"`
	Side           string            `ion:"side"`
	VenueOrderId   string            `ion:"venueOrderId"`
	FillId         string            `ion:"fillId"`
	Sender         *QldbCoreAccount  `ion:"sender"`
	Receiver       *QldbCoreAccount  `ion:"receiver"`
	FilledQuantity *ion.Decimal      `ion:"filledQuantity"`
	FilledValue    *ion.Decimal      `ion:"filledValue"`
	VenueFee       *ion.Decimal      `ion:"venueFee"`
	RetailFee      *ion.Decimal      `ion:"retailFee"`
	CreatedAt      time.Time         `ion:"createdAt"`
	Metadata       map[string]string `ion:"metadata"`
}

type QldbHold struct {
	HoldUUID   string       `ion:"holdUUID"`
	AccountId  string       `ion:"accountIndex1"`
	Amount     *ion.Decimal `ion:"amount"`
	ReleasedAt time.Time    `ion:"ReleasedAt"`
	Released   bool         `ion:"Released"`
}
