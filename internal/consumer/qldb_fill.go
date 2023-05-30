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
	"errors"
	"fmt"
	"time"

	"github.com/amzn/ion-go/ion"
	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/coinbase-samples/ib-ledger-go/internal/qldb"
	"github.com/coinbase-samples/ib-ledger-go/internal/relationaldb"
	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
	"github.com/google/uuid"
)

func writeQldbFill(ctx context.Context, f *model.QldbFill) error {
	if err := validateEntryIds(f.Metadata); err != nil {
		return fmt.Errorf(
			"invalid entry ids provided - fillId: %s - %w",
			f.FillId,
			err,
		)
	}

	qt, err := qldb.GetTransaction(ctx, f.VenueOrderId)
	if err != nil {
		return fmt.Errorf(
			"failed to get transaction - fillId: %s - orderId: %s - %w",
			f.FillId,
			f.VenueOrderId,
			err,
		)
	}

	sender, err := qt.Sender.ConvertToPostgresAccount()
	if err != nil {
		return fmt.Errorf(
			"bad sender account - fillId: %s - senderId: %s - %w",
			f.FillId,
			qt.Sender.Id,
			err,
		)
	}
	receiver, err := qt.Receiver.ConvertToPostgresAccount()
	if err != nil {
		return fmt.Errorf(
			"bad receiver account in fill - fill id: %s - receiver: %v - %w",
			f.Id,
			f.Receiver,
			err,
		)
	}
	transaction, err := qt.ConvertToPostgresTransaction()
	if err != nil {
		return fmt.Errorf(
			"bad transaction for fill - fill id: %s - transaction: %v - %w",
			f.Id,
			qt,
			err,
		)
	}

	if err := relationaldb.InsertTransaction(
		ctx,
		transaction,
		sender,
		receiver,
	); err != nil {
		return fmt.Errorf(
			"failed to upsert transaction - id: %s - %w",
			transaction.Id,
			err,
		)
	}

	fillUUID, err := uuid.Parse(f.FillId)
	if err != nil {
		return fmt.Errorf(
			"unable to parse fill id into valid uuid - fillId: %s - %w",
			f.FillId,
			err,
		)
	}

	var senderAmount *ion.Decimal
	var receiverAmount *ion.Decimal
	var venueFeeSenderAccountId uuid.UUID
	if f.Side == config.Buy {
		senderAmount = f.FilledValue
		receiverAmount = f.FilledQuantity
		venueFeeSenderAccountId = sender.Id
	} else {
		senderAmount = f.FilledQuantity
		receiverAmount = f.FilledValue
		venueFeeSenderAccountId = receiver.Id
	}

	senderEntryId := f.Metadata[config.SenderEntryKey]
	receiverEntryId := f.Metadata[config.ReceiverEntryKey]
	if err := insertEntries(
		ctx,
		senderEntryId,
		receiverEntryId,
		fillUUID,
		transaction.Id,
		sender.Id,
		receiver.Id,
		senderAmount,
		receiverAmount,
		f.CreatedAt,
	); err != nil {
		return fmt.Errorf(
			"failed to insert entries - fillId: %s - transactionId: %s - %w",
			f.FillId,
			transaction.Id.String(),
			err,
		)
	}

	// If VenueFee was taken out, write the corresponding entries
	if f.VenueFee != nil && f.VenueFee.Cmp(ion.NewDecimalInt(0)) == 1 {
		venueFeeSenderEntry := f.Metadata[config.VenueFeeSenderEntryKey]
		venueFeeReceiverEntry := f.Metadata[config.VenueFeeReceiverEntryKey]
		venueFeeReceiverAccountId, err := uuid.Parse(Repo.App.CoinbaseUsdAccountId)
		if err != nil {
			return fmt.Errorf(
				"unable to parse venueFee account string - %s",
				Repo.App.CoinbaseUsdAccountId,
			)
		}
		if err := insertEntries(
			ctx,
			venueFeeSenderEntry,
			venueFeeReceiverEntry,
			fillUUID,
			transaction.Id,
			venueFeeSenderAccountId,
			venueFeeReceiverAccountId,
			f.VenueFee,
			f.VenueFee,
			f.CreatedAt,
		); err != nil {
			return fmt.Errorf(
				"failed to insert venueFee entries - fillId: %s - transactionId: %s - %w",
				f.FillId,
				transaction.Id.String(),
				err,
			)
		}
	}

	// If Fee was taken out, write the corresponding entries
	if f.RetailFee != nil && f.RetailFee.Cmp(ion.NewDecimalInt(0)) == 1 {
		retailFeeSenderEntry := f.Metadata[config.RetailFeeSenderEntryKey]
		retailFeeReceiverEntry := f.Metadata[config.RetailFeeReceiverEntryKey]
		retailFeeReceiverAccountId, err := uuid.Parse(Repo.App.NeoworksUsdAccountId)
		if err != nil {
			return fmt.Errorf(
				"unable to parse fee account string - %s",
				Repo.App.NeoworksUsdAccountId,
			)
		}
		if err := insertEntries(
			ctx,
			retailFeeSenderEntry,
			retailFeeReceiverEntry,
			fillUUID,
			transaction.Id,
			venueFeeSenderAccountId,
			retailFeeReceiverAccountId,
			f.RetailFee,
			f.RetailFee,
			f.CreatedAt,
		); err != nil {
			return fmt.Errorf(
				"failed to insert fee entries - fillId: %s - transactionId - %s - %w",
				f.FillId,
				transaction.Id.String(),
				err,
			)
		}
	}

	return nil
}

func validateEntryIds(keyMap map[string]string) error {
	if _, ok := keyMap[config.SenderEntryKey]; !ok {
		return errors.New("no sender entry id")
	}
	if _, ok := keyMap[config.ReceiverEntryKey]; !ok {
		return errors.New("no receiver entry id")
	}
	_, venueFeeSenderOk := keyMap[config.VenueFeeSenderEntryKey]
	_, venueFeeReceiverOk := keyMap[config.VenueFeeReceiverEntryKey]
	if !venueFeeSenderOk && venueFeeReceiverOk {
		return errors.New("venueFee receiver entry present with no sender entry")
	}
	if venueFeeSenderOk && !venueFeeReceiverOk {
		return errors.New("venueFee sender entry present with no receiver entry")
	}

	_, feeSenderOk := keyMap[config.RetailFeeSenderEntryKey]
	_, feeReceiverOk := keyMap[config.RetailFeeReceiverEntryKey]
	if !feeSenderOk && feeReceiverOk {
		return errors.New("fee receiver entry present with no sender entry")
	}
	if feeSenderOk && !feeReceiverOk {
		return errors.New("fee sender entry present with no receiver entry")
	}
	return nil
}

func insertEntries(
	ctx context.Context,
	senderEntryId, receiverEntryId string,
	fillId, transactionId, senderAccountId, receiverAccountId uuid.UUID,
	senderAmount *ion.Decimal,
	receiverAmount *ion.Decimal,
	timestamp time.Time,
) error {
	receiverEntryUUID, err := uuid.Parse(receiverEntryId)
	if err != nil {
		return fmt.Errorf("unable to parse entry id into valid uuid: %w", err)
	}
	senderEntryUUID, err := uuid.Parse(senderEntryId)
	if err != nil {
		return fmt.Errorf("unable to parse entry id into valid uuid: %w", err)
	}
	senderAmountInt, err := utils.IonDecimalToBigInt(senderAmount)
	if err != nil {
		return fmt.Errorf(
			"unable to convert sender amount to int - amount: %s - %w",
			senderAmount.String(),
			err,
		)
	}
	receiverAmountInt, err := utils.IonDecimalToBigInt(receiverAmount)
	if err != nil {
		return fmt.Errorf(
			"unable to convert receiver amount to int - amount: %s - %w",
			receiverAmount.String(),
			err,
		)
	}

	if err := relationaldb.InsertEntry(
		ctx,
		&model.Entry{
			Id:           receiverEntryUUID,
			AccountId:    receiverAccountId,
			VenueOrderId: transactionId,
			Amount:       receiverAmountInt.String(),
			Direction:    config.Credit,
			CreatedAt:    timestamp,
			FillId:       fillId,
		},
		&model.Entry{
			Id:           senderEntryUUID,
			AccountId:    senderAccountId,
			VenueOrderId: transactionId,
			Amount:       senderAmountInt.String(),
			Direction:    config.Debit,
			CreatedAt:    timestamp,
			FillId:       fillId,
		},
	); err != nil {
		return fmt.Errorf("failed to insert entry to database: %w", err)
	}
	return nil
}
