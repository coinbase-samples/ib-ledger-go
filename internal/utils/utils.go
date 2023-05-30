/**
 * Copyright 2022-present Coinbase Global, Inc.
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
	"crypto/md5"
	"fmt"
	"math/big"
	"strings"

	"github.com/amzn/ion-go/ion"
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
)

func GetTransactionTypeFromString(s string) (api.TransactionType, bool) {
	switch s {
	case "TRANSFER":
		return api.TransactionType_TRANSACTION_TYPE_TRANSFER, true
	case "ORDER":
		return api.TransactionType_TRANSACTION_TYPE_ORDER, true
	case "CONVERT":
		return api.TransactionType_TRANSACTION_TYPE_CONVERT, true
	default:
		return api.TransactionType_TRANSACTION_TYPE_UNSPECIFIED, false
	}
}

func GetStringFromTransactionType(t api.TransactionType) (string, bool) {
	switch t {
	case api.TransactionType_TRANSACTION_TYPE_TRANSFER:
		return "TRANSFER", true
	case api.TransactionType_TRANSACTION_TYPE_ORDER:
		return "ORDER", true
	case api.TransactionType_TRANSACTION_TYPE_CONVERT:
		return "CONVERT", true
	default:
		return "", false
	}
}

func GenerateIdemString(params ...string) string {
	hashedBytes := md5.Sum([]byte(strings.Join(params, "")))
	return fmt.Sprintf("%x", hashedBytes)
}

func FormatProductId(debitCurrency, creditCurrency string) string {
	return fmt.Sprintf("%s-%s", strings.ToUpper(creditCurrency), strings.ToUpper(debitCurrency))
}

func SplitProductId(productId string) (string, string) {
	split := strings.Split(productId, "-")
	if len(split) != 2 {
		return "", ""
	}
	return split[0], split[1]
}

func IonDecimalToBigInt(i *ion.Decimal) (*big.Int, error) {
	out, exp := i.CoEx()
	if exp != 0 {
		return nil, fmt.Errorf(
			"bad input, ion decimal should only be wrapping an int: %s",
			i.String(),
		)
	}

	return out, nil
}

func SetString(i *big.Int, s string) error {
	if _, ok := i.SetString(s, 10); !ok {
		return fmt.Errorf("failed to convert %s to bigInt", s)
	}
	return nil
}
