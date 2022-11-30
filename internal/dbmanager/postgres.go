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

	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type PostgresDBManager struct {
	Pool *pgxpool.Pool
}

func NewPostgresDBManager(app config.AppConfig) *PostgresDBManager {

	if app.DbCreds == "" {
		log.Fatalf("no environment variable set for DB_CREDENTIALS")
	}

	if app.DbHostname == "" {
		log.Fatalf("no environment variable set for DB_HOSTNAME")
	}

	if app.DbPort == "" {
		log.Fatalf("no environment variable set for DB_PORT")
	}

	var dbCredsJson map[string]interface{}
	err := json.Unmarshal([]byte(app.DbCreds), &dbCredsJson)

	if err != nil {
		log.Fatalf("unable to unmarshal the cred string")
	}

	dbUsername := dbCredsJson["username"].(string)
	dbPassword := url.QueryEscape(dbCredsJson["password"].(string))

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/ledger", dbUsername, dbPassword, app.DbHostname, app.DbPort)

	if app.Env == "local" {
		dbUrl += "?sslmode=disable"
	}

	log.Printf("attempting to connect to database with username: %v, hostname: %v, and port %v", dbUsername, app.DbHostname, app.DbPort)
	pool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		log.Fatalf("Failed to establish DB Pool: %v", err)
	}

	return &PostgresDBManager{Pool: pool}
}

func (d *PostgresDBManager) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, d.Pool, dest, query, args...)
}
