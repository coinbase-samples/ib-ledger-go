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

package dbmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/coinbase-samples/ib-ledger-go/config"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type PostgresDBManager struct {
	Pool *pgxpool.Pool
}

func NewPostgresDBManager(app *config.AppConfig, l *log.Entry) *PostgresDBManager {

	var dbCredsJson map[string]interface{}
	err := json.Unmarshal([]byte(app.DbCreds), &dbCredsJson)
	if err != nil {
		l.Fatalf("unable to unmarshal the cred string")
	}

	dbUsername := dbCredsJson["username"].(string)
	dbPassword := url.QueryEscape(dbCredsJson["password"].(string))

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/ledger", dbUsername, dbPassword, app.DbHostname, app.DbPort)

	if app.IsLocalEnv() {
		dbUrl += "?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		l.Fatalf("failed to establish DB Pool: %w", err)
	}

	return &PostgresDBManager{Pool: pool}
}

func (d *PostgresDBManager) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, d.Pool, dest, query, args...)
}
