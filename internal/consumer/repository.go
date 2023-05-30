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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	log "github.com/sirupsen/logrus"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	Svc *sqs.Client
}

func NewRepo(a *config.AppConfig, cfg *aws.Config) {
	Repo = &Repository{
		App: a,
		Svc: sqs.NewFromConfig(*cfg),
	}
}

func (r Repository) DeleteSqsMsg(
	ctx context.Context,
	receipt string,
) error {

	if _, err := r.Svc.DeleteMessage(
		ctx,
		&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(r.App.QueueUrl),
			ReceiptHandle: aws.String(receipt),
		},
	); err != nil {
		return fmt.Errorf("unable to delete msg: %w", err)
	}
	return nil
}

func (r Repository) ReceiveSqsMsg(
	ctx context.Context,
	timeout int32,
) (*types.Message, error) {
	log.Debugf("receiveSqs - %s - %v", r.App.QueueUrl, r.Svc)
	res, err := r.Svc.ReceiveMessage(
		ctx,
		&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(r.App.QueueUrl),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     0,
			VisibilityTimeout:   timeout,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("unable to receive sqs message: %w", err)
	}

	if len(res.Messages) == 0 {
		return nil, nil
	}

	return &(res.Messages[0]), nil
}

type Event struct {
	// The data blob. The data in the blob is both opaque and immutable to Kinesis
	// Data Streams, which does not inspect, interpret, or change the data in the
	// blob in any way. When the data blob (the payload before base64-encoding)
	// is added to the partition key size, the total size must not exceed the maximum
	// record size (1 MiB).
	// Data is automatically base64 encoded/decoded by the SDK.
	//
	// Data is a required field
	Data string `json:"data"`
}
