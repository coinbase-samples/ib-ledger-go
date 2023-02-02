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

package main

import (
	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	"github.com/coinbase-samples/ib-ledger-go/internal/consumer"
	"github.com/coinbase-samples/ib-ledger-go/internal/qldb"
	"github.com/coinbase-samples/ib-ledger-go/internal/relationaldb"
)

func main() {

	var app config.AppConfig

	config.Setup(&app)
	logrusLogger := config.LogInit(app)

	config.ValidateConfig(app, logrusLogger)

	cfg := app.GenerateAwsConfig(logrusLogger)

	qldb.NewRepository(logrusLogger, &app, &cfg)
	relationaldb.NewRepo(&app, logrusLogger)
	consumer.NewRepo(&app, &cfg)

	go consumer.Listen(logrusLogger)

	gRPCListen(app, logrusLogger)
}
