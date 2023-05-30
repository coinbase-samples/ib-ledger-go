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

import "github.com/amzn/ion-go/ion"

type StreamRecord struct {
	QldbStreamArn string                 `ion:"qldbStreamArn"`
	RecordType    string                 `ion:"recordType"`
	Payload       map[string]interface{} `ion:"payload"`
}

type TableInfo struct {
	TableName string `ion:"tableName"`
	TableId   string `ion:"tableId"`
}

type RevisionBlockAddress struct {
	StrandId   string `ion:"strandId"`
	SequenceNo uint64 `ion:"sequenceNo"`
}

type RevisionMetadata struct {
	Id      string        `ion:"id"`
	Version uint64        `ion:"version"`
	TxnTime ion.Timestamp `ion:"txnTime"`
	TxnId   string        `ion:"txId"`
}

type Revision struct {
	BlockAddress     RevisionBlockAddress `ion:"blockAddress"`
	Hash             []byte               `ion:"hash"`
	Data             interface{}          `ion:"data"`
	RevisionMetadata RevisionMetadata     `ion:"metadata"`
}

type RevisionRecordPayload struct {
	TableInfo TableInfo `ion:"tableInfo"`
	Revision  Revision  `ion:"revision"`
}
