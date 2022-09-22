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

package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/coinbase-samples/ib-ledger-go/config"
	"github.com/coinbase-samples/ib-ledger-go/model"
	"github.com/coinbase-samples/ib-ledger-go/utils"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	api "github.com/coinbase-samples/ib-ledger-go/protos/ledger"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func NewPostgresHandler(app config.AppConfig) *PostgresRepository {

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
	pool, err := pgxpool.Connect(context.Background(), dbUrl)

	if err != nil {
		log.Fatalf("Failed to establish DB Pool: %v", err)
	}

	return &PostgresRepository{Pool: pool}
}

func (handler *PostgresRepository) InitializeAccount(ctx context.Context, request *api.InitializeAccountRequest) (*model.InitializeAccountResult, error) {
	var initializeAccountResult []*model.InitializeAccountResult
	const sql = `SELECT id, portfolio_id, user_id, currency, created_at, balance, hold, available FROM initialize_account($1, $2, $3)`

	err := pgxscan.Select(context.Background(), handler.Pool, &initializeAccountResult, sql, request.PortfolioId, request.UserId, request.Currency)
	if err != nil {
		return nil, err
	}

	return initializeAccountResult[0], nil
}

func (handler *PostgresRepository) CreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (*model.CreateTransactionResult, error) {
	var createTransactionResult []*model.CreateTransactionResult

	const sql = `SELECT id, sender_id, receiver_id, request_id, transaction_type, created_at FROM create_transaction_and_place_hold($1, $2, $3, $4, $5, $6, $7, $8)`

	transactionType, ok := utils.GetStringFromTransactionType(request.TransactionType)
	if !ok {
		return nil, fmt.Errorf("bad request: transaction type not supported: %v", transactionType)
	}

	err := pgxscan.Select(context.Background(),
		handler.Pool,
		&createTransactionResult,
		sql,
		request.OrderId,
		request.Sender.Currency,
		request.Sender.UserId,
		request.Receiver.Currency,
		request.Receiver.UserId,
		request.RequestId.Value,
		request.TotalAmount,
		transactionType)
	if err != nil {
		return nil, err
	}

	return createTransactionResult[0], nil
}

func (handler *PostgresRepository) PartialReleaseHold(ctx context.Context, request *api.PartialReleaseHoldRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	const sql = `SELECT hold_id, sender_entry_id, receiver_entry_id, sender_balance_id, receiver_balance_id FROM partial_release_hold($1, $2, $3, $4)`

	err := pgxscan.Select(context.Background(), handler.Pool, &TransactionResult, sql, request.OrderId, request.RequestId, request.SenderAmount, request.ReceiverAmount)
	if err != nil {
		return nil, err
	}

	return TransactionResult[0], nil
}

func (handler *PostgresRepository) CompleteTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	const sql = `SELECT hold_id, sender_entry_id, receiver_entry_id, sender_balance_id, receiver_balance_id FROM complete_transaction($1, $2, $3)`

	err := pgxscan.Select(context.Background(), handler.Pool, &TransactionResult, sql, request.OrderId, request.RequestId, request.ReceiverAmount)
	if err != nil {
		return nil, err
	}

	return TransactionResult[0], nil
}

func (handler *PostgresRepository) FailTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	const sql = `SELECT hold_id, sender_entry_id, receiver_entry_id, sender_balance_id, receiver_balance_id FROM fail_transaction($1, $2)`

	err := pgxscan.Select(context.Background(), handler.Pool, &TransactionResult, sql, request.OrderId, request.RequestId)
	if err != nil {
		return nil, err
	}

	return TransactionResult[0], nil
}

func (handler *PostgresRepository) CancelTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	const sql = `SELECT hold_id, sender_entry_id, receiver_entry_id, sender_balance_id, receiver_balance_id FROM cancel_transaction($1, $2)`

	err := pgxscan.Select(context.Background(), handler.Pool, &TransactionResult, sql, request.OrderId, request.RequestId)
	if err != nil {
		return nil, err
	}

	return TransactionResult[0], nil
}

func (handler *PostgresRepository) GetAllAccountsAndMostRecentBalances(ctx context.Context, userId string) ([]*model.GetAccountResult, error) {
	var accountResult []*model.GetAccountResult
	l := ctxlogrus.Extract(ctx)

	const sql = `SELECT account_id, currency, balance, hold, available, created_at FROM get_balances_for_users($1)`

	err := pgxscan.Select(context.Background(), handler.Pool, &accountResult, sql, userId)
	if err != nil {
		l.Debugln("error getting accounts", err)
		return nil, err
	}
	l.Debugln("fetched accounts and balances", accountResult)
	return accountResult, nil
}
