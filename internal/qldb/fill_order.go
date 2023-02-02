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
	"strings"
	"time"

	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/google/uuid"

	"github.com/amzn/ion-go/ion"
	"github.com/awslabs/amazon-qldb-driver-go/v3/qldbdriver"
	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
)

func FillOrder(req *api.PostFillRequest) error {
	_, err := Repo.Driver.Execute(
		context.Background(),
		func(txn qldbdriver.Transaction) (interface{}, error) {
			t, s, r, err := getTransactionAndAccountsForFill(
				txn,
				req.OrderId,
			)
			if err != nil {
				return nil, err
			}

			// check for fill
			fillId := strings.ToLower(req.FillId)
			fillQldbId := model.GenerateFillId(t.VenueOrderId, fillId)
			result, err := txn.Execute(
				"SELECT * FROM Ledger WHERE id = ?",
				fillQldbId,
			)

			if err != nil {
				return nil, fmt.Errorf(
					"failed to query ledger for fill - fillId: %s - transactionId: %s - err: %w",
					fillId,
					t.VenueOrderId,
					err,
				)
			}

			if result.Next(txn) {
				return nil, nil
			}

			filledValue, err := ion.ParseDecimal(req.FilledValue)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to parse filledValue into ion decimal - filledValue: %s - err: %w",
					req.FilledValue,
					err,
				)
			}
			filledQuantity, err := ion.ParseDecimal(req.FilledQuantity)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to parse filledQuantity into ion decimal - filledQuantity: %s - err: %w",
					req.FilledQuantity,
					err,
				)
			}

			metadata := make(map[string]string)
			venueFee, _ := ion.ParseDecimal("0")
			if req.VenueFeeAmount != nil {
				venueFee, err = handleFeeFills(
					txn,
					req.VenueFeeAmount.Value,
					Repo.App.CoinbaseUserId,
				)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to process fill for venue fee: %w",
						err,
					)
				}
				metadata[config.VenueFeeSenderEntryKey] = uuid.New().String()
				metadata[config.VenueFeeReceiverEntryKey] = uuid.New().String()
			}

			retailFee, _ := ion.ParseDecimal("0")
			if req.RetailFeeAmount != nil {
				retailFee, err = handleFeeFills(
					txn,
					req.RetailFeeAmount.Value,
					Repo.App.NeoworksUserId,
				)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to process fill for retail fee: %w",
						err,
					)
				}
				metadata[config.RetailFeeReceiverEntryKey] = uuid.New().String()
				metadata[config.RetailFeeSenderEntryKey] = uuid.New().String()
			}

			metadata[config.SenderEntryKey] = uuid.New().String()
			metadata[config.ReceiverEntryKey] = uuid.New().String()
			fill := &model.QldbFill{
				Id:             fillQldbId,
				ProductId:      req.ProductId,
				Side:           req.Side,
				VenueOrderId:   t.VenueOrderId,
				FillId:         fillId,
				Sender:         s.GetCoreAccount(),
				Receiver:       r.GetCoreAccount(),
				FilledQuantity: filledQuantity,
				FilledValue:    filledValue,
				VenueFee:       venueFee,
				RetailFee:      retailFee,
				CreatedAt:      time.Now(),
				Metadata:       metadata,
			}

			if _, err := txn.Execute("INSERT INTO Ledger ?", fill); err != nil {
				return nil, fmt.Errorf(
					"failed to insert fill into ledger - fill: %v - error: %w",
					fill,
					err,
				)
			}

			if err := processFillAccountAndHoldUpdates(
				txn,
				req.Side,
				t,
				s,
				r,
				filledQuantity,
				filledValue,
				venueFee,
				retailFee,
			); err != nil {
				return nil, fmt.Errorf(
					"failed to process account balance and hold update: %w",
					err,
				)
			}

			return nil, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed qldb fill operations - fillId: %s - error: %w",
			req.FillId,
			err,
		)
	}
	return nil
}

func getTransactionAndAccountsForFill(
	txn qldbdriver.Transaction,
	orderId string,
) (*model.QldbTransaction, *model.QldbAccount, *model.QldbAccount, error) {
	venueOrderId := strings.ToLower(orderId)
	t, err := getTransactionQldbOperation(txn, venueOrderId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get transaction from fill - %w",
			err,
		)
	}

	s, err := getAccountQldbTransaction(txn, t.Sender.Id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get sender account - %w",
			err,
		)
	}

	r, err := getAccountQldbTransaction(txn, t.Receiver.Id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get receiver account - %w",
			err,
		)
	}
	return t, s, r, nil
}

func handleFeeFills(
	txn qldbdriver.Transaction,
	amount, userId string,
) (*ion.Decimal, error) {
	fee, err := ion.ParseDecimal(amount)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse fee into ion decimal - fee: %s - err: %w",
			amount,
			err,
		)
	}
	feeAccountId := model.GenerateAccountId(
		userId,
		"USD",
	)
	f, err := getAccountQldbTransaction(txn, feeAccountId)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get fee account id - %w",
			err,
		)
	}

	feeInt, err := utils.IonDecimalToBigInt(fee)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to convert ion decimal to int - %w",
			err,
		)
	}

	if err := creditAccountUpdate(txn, f, feeInt); err != nil {
		return nil, fmt.Errorf(
			"failed to update account balance for fee - accountId: - %s - error: %w",
			f.Id,
			err,
		)

	}
	return fee, nil
}

func processFillAccountAndHoldUpdates(
	txn qldbdriver.Transaction,
	side string,
	t *model.QldbTransaction,
	s, r *model.QldbAccount,
	filledQuantity, filledValue, venueFee, retailFee *ion.Decimal,
) error {

	filledValueInt, err := utils.IonDecimalToBigInt(filledValue)
	if err != nil {
		return err
	}

	filledQuantityInt, err := utils.IonDecimalToBigInt(filledQuantity)
	if err != nil {
		return err
	}

	venueFeeInt, _ := utils.IonDecimalToBigInt(venueFee)
	if err != nil {
		return err
	}

	retailFeeInt, err := utils.IonDecimalToBigInt(retailFee)
	if err != nil {
		return err
	}

	if side == config.Buy {
		if err := processBuySideAccountUpdates(
			txn,
			s,
			r,
			filledQuantityInt,
			filledValueInt,
			venueFeeInt,
			retailFeeInt,
		); err != nil {
			return fmt.Errorf(
				"failed to process buy side account update: %w",
				err,
			)
		}
	} else {
		if err := processSellSideAccountUpdates(
			txn,
			s,
			r,
			filledQuantityInt,
			filledValueInt,
			venueFeeInt,
			retailFeeInt,
		); err != nil {
			return fmt.Errorf(
				"failed to process sell side account updates: %w",
				err,
			)
		}
	}

	// Update Hold
	hAmount, err := utils.IonDecimalToBigInt(t.Hold.Amount)
	if err != nil {
		return err
	}

	hAmount.Sub(hAmount, filledValueInt)
	hAmount.Sub(hAmount, venueFeeInt)
	hAmount.Sub(hAmount, retailFeeInt)

	hold := t.Hold
	hold.Amount = ion.NewDecimal(hAmount, 0, false)

	if _, err := txn.Execute(
		"UPDATE Ledger AS t SET t.hold = ? WHERE t.id = ?",
		hold,
		t.Id,
	); err != nil {

		return fmt.Errorf(
			"failed to update hold for fill - transactionId: - %s - error: %w",
			t.Id,
			err,
		)
	}
	return nil
}

func processBuySideAccountUpdates(
	txn qldbdriver.Transaction,
	s, r *model.QldbAccount,
	filledQuantity, filledValue, venueFee, retailFee *big.Int,
) error {
	if err := debitAccountUpdate(
		txn,
		s,
		filledValue,
		venueFee,
		retailFee,
	); err != nil {
		return fmt.Errorf(
			"failed sender account balance update - error: %w",
			err,
		)
	}

	if err := creditAccountUpdate(
		txn,
		r,
		filledQuantity,
	); err != nil {
		return fmt.Errorf(
			"failed receiver account balance update - error: %w",
			err,
		)
	}
	return nil
}

func processSellSideAccountUpdates(
	txn qldbdriver.Transaction,
	s, r *model.QldbAccount,
	filledQuantity, filledValue, venueFee, retailFee *big.Int,
) error {
	// when selling assets, fees are paid in the denominator (USD)
	if err := debitAccountUpdate(
		txn,
		s,
		filledQuantity,
		big.NewInt(0),
		big.NewInt(0),
	); err != nil {
		return fmt.Errorf(
			"failed sender account balance update - error: %w",
			err,
		)
	}

	if err := creditAccountUpdate(
		txn,
		r,
		filledValue,
	); err != nil {
		return fmt.Errorf(
			"failed receiver account balance update - error: %w",
			err,
		)
	}

	if err := debitAccountUpdate(
		txn,
		r,
		big.NewInt(0),
		venueFee,
		retailFee,
	); err != nil {
		return fmt.Errorf(
			"failed receiver fee account balance update - error: %w",
			err,
		)
	}
	return nil
}
