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

package utils

import (
	api "LedgerApp/protos/ledger"
)

func GetTransactionTypeFromString(s string) (api.TransactionType, bool) {
	switch s {
	case "TRANSFER":
		return api.TransactionType_TRANSFER, true
	case "ORDER":
		return api.TransactionType_ORDER, true
	case "CONVERT":
		return api.TransactionType_CONVERT, true
	default:
		return api.TransactionType_TRANSFER, false
	}
}

func GetStringFromTransactionType(t api.TransactionType) (string, bool) {
	switch t {
	case api.TransactionType_TRANSFER:
		return "TRANSFER", true
	case api.TransactionType_ORDER:
		return "ORDER", true
	case api.TransactionType_CONVERT:
		return "CONVERT", true
	default:
		return "", false
	}
}
