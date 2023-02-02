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

package relationaldb

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

var Repo *Repository

type Repository struct {
	Pool *pgxpool.Pool
}

type DbManager interface {
	Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Insert(ctx context.Context, sql string, args ...interface{}) error
}

func NewRepo(a *config.AppConfig, l *log.Entry) {
	var dbCredsJson map[string]interface{}
	err := json.Unmarshal([]byte(a.DbCreds), &dbCredsJson)
	if err != nil {
		l.Fatal("unable to unmarshal the database credentials string")
	}

	dbUsername := dbCredsJson["username"].(string)
	dbPassword := url.QueryEscape(dbCredsJson["password"].(string))

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/ledger",
		dbUsername,
		dbPassword,
		a.DbHostname,
		a.DbPort,
	)

	if a.IsLocalEnv() {
		dbUrl += "?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		l.Fatalf("failed to establish DB Pool: %v", err.Error())
	}

	Repo = &Repository{
		Pool: pool,
	}
}

func (r *Repository) Query(
	ctx context.Context,
	dest interface{},
	query string,
	args ...interface{},
) error {
	return pgxscan.Select(ctx, r.Pool, dest, query, args...)
}

func (r *Repository) Insert(
	ctx context.Context,
	sql string,
	args ...interface{},
) error {
	_, err := r.Pool.Exec(ctx, sql, args...)
	return err
}
