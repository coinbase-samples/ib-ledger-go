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
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/amzn/ion-go/ion"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/awslabs/kinesis-aggregation/go/v2/deaggregator"
	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/coinbase-samples/ib-ledger-go/internal/relationaldb"

	log "github.com/sirupsen/logrus"
)

func Listen(l *log.Entry) {
	for {
		ctx := context.Background()
		message, err := Repo.ReceiveSqsMsg(ctx, 90)

		if err != nil {
			time.Sleep(500 * time.Millisecond)
			l.Error(err)
			continue
		}

		if message == nil {
			continue
		}

		if err := onMessage(
			ctx,
			message.Body,
		); err != nil {
			time.Sleep(500 * time.Millisecond)
			l.Error(err)
			continue
		}

		if err := Repo.DeleteSqsMsg(ctx, *message.ReceiptHandle); err != nil {
			time.Sleep(500 * time.Millisecond)
			l.Error(err)
			continue
		}

	}
}

func onMessage(ctx context.Context, body *string) error {
	var kinesisEvent events.KinesisRecord
	if err := json.Unmarshal([]byte(*body), &kinesisEvent); err != nil {
		return fmt.Errorf("unable to parse sqs into kinesis event: %w", err)
	}

	kr := &types.Record{
		ApproximateArrivalTimestamp: &kinesisEvent.ApproximateArrivalTimestamp.Time,
		Data:                        kinesisEvent.Data,
		EncryptionType:              getKinesisEncryptionType(kinesisEvent),
		PartitionKey:                &kinesisEvent.PartitionKey,
		SequenceNumber:              &kinesisEvent.SequenceNumber,
	}

	kcr := []types.Record{*kr}

	dars, err := deaggregator.DeaggregateRecords(kcr)
	if err != nil {
		return fmt.Errorf("failed to deaggregate records: %w", err)
	}

	for _, r := range dars {
		if err := handleRecordData(ctx, r.Data); err != nil {
			return err
		}
	}
	return nil
}

func getKinesisEncryptionType(
	kinesisEvent events.KinesisRecord,
) types.EncryptionType {
	if kinesisEvent.EncryptionType == "KMS" {
		return types.EncryptionTypeKms
	}
	return types.EncryptionTypeNone
}

func getLedgerType(reader ion.Reader) (string, error) {
	for reader.Next() {
		annotations, err := reader.Annotations()
		if err != nil {
			return "", fmt.Errorf(
				"unable to get annotations from ion blob: %w",
				err,
			)
		}
		if annotations == nil && reader.Type() == ion.StructType {
			reader.StepIn()
			for reader.Next() {
				fieldName, err := reader.FieldName()
				if err != nil {
					return "", fmt.Errorf(
						"unable to extract field name from ion node in ledger blob: %w",
						err,
					)
				}
				if *fieldName.Text == "id" {
					id, err := reader.StringValue()
					if err != nil {
						return "", fmt.Errorf(
							"unable to parse string from Id field: %w",
							err,
						)
					}
					return *id, nil
				}
			}
			reader.StepOut()
		}
	}
	return "", errors.New("no Id field found in ion blob")
}

func handleRecordData(ctx context.Context, blob []byte) error {
	var data map[string]interface{}
	if err := ion.Unmarshal(blob, &data); err != nil {
		return fmt.Errorf("failed to unmarshal ion data: %w", err)
	}

	if data["recordType"] != "REVISION_DETAILS" {
		return nil
	}

	payloadBytes, err := ion.MarshalBinary(data["payload"])
	if err != nil {
		return fmt.Errorf("unable to marshal payload into bytes: %w", err)
	}

	var payload *RevisionRecordPayload
	if err := ion.Unmarshal(payloadBytes, &payload); err != nil {
		return fmt.Errorf("unable to unmarshal payload into struct: %w", err)
	}

	payloadData, err := ion.MarshalBinary(payload.Revision.Data)
	if err != nil {
		return fmt.Errorf("unable to marshal payload data: %w", err)
	}

	switch payload.TableInfo.TableName {
	case config.AccountsTable:
		return handleQldbAccount(ctx, payloadData)
	case config.LedgerTable:
		return handleLedgerTableTypes(ctx, payloadData)
	default:
		return fmt.Errorf("bad table name: %s", payload.TableInfo.TableName)
	}
}

func handleQldbAccount(ctx context.Context, data []byte) error {
	var a *model.QldbAccount
	if err := ion.Unmarshal(data, &a); err != nil {
		return fmt.Errorf(
			"unable to unmarshal account blob into struct: %w",
			err,
		)
	}
	if err := relationaldb.InsertAccountBalance(ctx, a); err != nil {
		return fmt.Errorf(
			"unable to insert account balance into database: %w",
			err,
		)
	}
	return nil
}

func handleLedgerTableTypes(ctx context.Context, data []byte) error {
	reader := ion.NewReaderBytes(data)
	ledgerType, err := getLedgerType(reader)
	if err != nil {
		return fmt.Errorf("unable to get ledgerType from blob: %w", err)
	}
	splitId := strings.Split(ledgerType, "#")
	prefix := splitId[0]
	switch prefix {
	case "transaction":
		return handleQldbTransaction(ctx, data)
	case "fill":
		return handleQldbFill(ctx, data)
	default:
		return fmt.Errorf("bad prefix for ledger id: %s", prefix)
	}
}

func handleQldbTransaction(ctx context.Context, data []byte) error {
	var t *model.QldbTransaction
	if err := ion.Unmarshal(data, &t); err != nil {
		return fmt.Errorf(
			"unable to unmarshal transaction from data: %w",
			err,
		)
	}
	if err := writeQldbTransaction(ctx, t); err != nil {
		return fmt.Errorf(
			"unable to write Transaction to database: %w",
			err,
		)
	}
	return nil
}

func handleQldbFill(ctx context.Context, data []byte) error {
	var f *model.QldbFill
	if err := ion.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("unable to unmarshal fill from data: %w", err)
	}
	if err := writeQldbFill(ctx, f); err != nil {
		return fmt.Errorf("unable to write Fill data to database: %w", err)
	}
	return nil
}
