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
	"fmt"

	ledgererr "github.com/coinbase-samples/ib-ledger-go/internal/errors"
	"google.golang.org/grpc/codes"
)

const NeoworksUsdAccount = "B72D0E55-F53A-4DB0-897E-2CE4A73CB94B"

const CoinbaseUsdAccount = "C4D0E14E-1B2B-4023-AFA6-8891AD1960C9"

func GetFeeAccounts(currency string) (string, string, error) {
	switch currency {
	case "USD":
		return NeoworksUsdAccount, CoinbaseUsdAccount, nil
	default:
		return "", "", ledgererr.New(codes.InvalidArgument, fmt.Sprintf("invalid currency for fee accounts: %v", currency))
	}
}
