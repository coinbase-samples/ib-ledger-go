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

package main

import (
	"context"
	"time"

	"github.com/amzn/ion-go/ion"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/qldbsession"
	"github.com/awslabs/amazon-qldb-driver-go/v3/qldbdriver"
	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	log "github.com/sirupsen/logrus"
)

func main() {
	var app config.AppConfig

	config.Setup(&app)
	logrusLogger := config.LogInit(app)

	cfg := app.GenerateAwsConfig(logrusLogger)

	driver := initializeQldbDriver(logrusLogger, app, cfg)
	defer driver.Shutdown(context.Background())

	initializeLedgerTable(logrusLogger, driver)

	initializeAccountTableAndFeeAccounts(logrusLogger, driver, &app)
}

func initializeQldbDriver(
	l *log.Entry,
	app config.AppConfig,
	cfg aws.Config,
) *qldbdriver.QLDBDriver {
	qldbSession := qldbsession.NewFromConfig(cfg, func(options *qldbsession.Options) {
		options.Region = app.DevRegion
	})
	driver, err := qldbdriver.New(
		app.QldbName,
		qldbSession,
		func(options *qldbdriver.DriverOptions) {
			options.LoggerVerbosity = qldbdriver.LogInfo
		},
	)
	if err != nil {
		l.Fatalf("failed to create qldbdriver: %v", err)
	}
	return driver
}

func initializeLedgerTable(l *log.Entry, driver *qldbdriver.QLDBDriver) {
	if _, err := driver.Execute(
		context.Background(),
		func(txn qldbdriver.Transaction) (interface{}, error) {
			if _, err := txn.Execute(
				"CREATE TABLE Ledger",
			); err != nil {
				return nil, err
			}

			if _, err := txn.Execute(
				"CREATE INDEX ON Ledger (id)",
			); err != nil {
				return nil, err
			}

			if _, err := txn.Execute(
				"CREATE INDEX ON Ledger (venueOrderId)",
			); err != nil {
				return nil, err
			}

			if _, err := txn.Execute(
				"CREATE INDEX ON Ledger (fillId)",
			); err != nil {
				return nil, err
			}

			if _, err := txn.Execute(
				"CREATE INDEX ON Ledger (debitAccount)",
			); err != nil {
				return nil, err
			}

			_, err := txn.Execute("CREATE INDEX ON Ledger (creditAccount)")
			return nil, err
		},
	); err != nil {
		l.Fatalf("failed to create Ledger table: %v", err)
	}
}

func initializeAccountTableAndFeeAccounts(
	l *log.Entry,
	driver *qldbdriver.QLDBDriver,
	app *config.AppConfig,
) {
	coinbaseUsdFeeAccount := &model.QldbAccount{
		Id:          model.GenerateAccountId(app.CoinbaseUserId, "USD"),
		UserId:      app.CoinbaseUserId,
		Currency:    "USD",
		Balance:     ion.NewDecimalInt(0),
		Hold:        ion.NewDecimalInt(0),
		Available:   ion.NewDecimalInt(0),
		AccountUUID: app.CoinbaseUsdAccountId,
		UpdatedAt:   time.Now(),
	}

	neoworksUsdFeeAccount := &model.QldbAccount{
		Id:          model.GenerateAccountId(app.NeoworksUserId, "USD"),
		UserId:      app.NeoworksUserId,
		Currency:    "USD",
		Balance:     ion.NewDecimalInt(0),
		Hold:        ion.NewDecimalInt(0),
		Available:   ion.NewDecimalInt(0),
		AccountUUID: app.NeoworksUsdAccountId,
		UpdatedAt:   time.Now(),
	}

	if _, err := driver.Execute(
		context.Background(),
		func(
			txn qldbdriver.Transaction,
		) (interface{}, error) {
			if _, err := txn.Execute("CREATE TABLE Accounts"); err != nil {
				return nil, err
			}

			if _, err := txn.Execute(
				"CREATE INDEX ON Accounts (id)",
			); err != nil {
				return nil, err
			}

			if _, err := txn.Execute(
				"CREATE INDEX ON Accounts (userId)",
			); err != nil {
				return nil, err
			}

			if _, err := txn.Execute(
				"INSERT INTO Accounts ?",
				[]*model.QldbAccount{
					coinbaseUsdFeeAccount,
					neoworksUsdFeeAccount,
				},
			); err != nil {
				return nil, err
			}

			return nil, nil
		},
	); err != nil {
		l.Fatalf("failed to create Accounts table: %v", err)
	}
}
