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
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	ledgererr "github.com/coinbase-samples/ib-ledger-go/internal/errors"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
	"google.golang.org/grpc/codes"

	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

var (
	initializeAccountSql = `
    SELECT id, portfolio_id, user_id, currency, created_at, balance, hold, available 
    FROM initialize_account($1, $2, $3);`

	getAllAccountsAndMostRecentBalancesSql = `
    SELECT account_id, currency, balance, hold, available, created_at 
    FROM get_balances_for_users($1);`

	createTransactionSql = `
    SELECT id, sender_id, receiver_id, request_id, transaction_type, created_at 
    FROM create_transaction_and_place_hold($1, $2, $3, $4, $5, $6, $7, $8);`

	partialReleaseHoldSql = `
    SELECT hold_id, sender_entry_id, receiver_entry_id, sender_balance_id, receiver_balance_id 
    FROM partial_release_hold($1, $2, $3, $4, $5, $6, $7, $8);`

	completeTransactionSql = `
    SELECT hold_id, sender_entry_id, receiver_entry_id, sender_balance_id, receiver_balance_id 
    FROM complete_transaction($1, $2);`

	failTransactionSql = `
    SELECT hold_id, sender_entry_id, receiver_entry_id, sender_balance_id, receiver_balance_id 
    FROM fail_transaction($1, $2);`

	cancelTransactionSql = `
    SELECT hold_id, sender_entry_id, receiver_entry_id, sender_balance_id, receiver_balance_id 
    FROM cancel_transaction($1, $2);`
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
	pool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		log.Fatalf("Failed to establish DB Pool: %v", err)
	}

	return &PostgresRepository{Pool: pool}
}

func (r *PostgresRepository) handleErrors(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		msg := pgErr.Message
		log.Errorf("postgres error - code: %v - message: %v", pgErr.SQLState(), msg)
        switch msg {
            case "insufficient available balance": 
			    return ledgererr.New(codes.InvalidArgument, "insufficient balance for transaction")
            case "sender account missing":
            case "receiver account missing": 
			    return ledgererr.New(codes.InvalidArgument, "account not found")
            default:
			    return ledgererr.New(codes.Internal, "postgres internal failure")
		}
	}
	return ledgererr.FromError(err)
}

func (handler *PostgresRepository) InitializeAccount(ctx context.Context, request *api.InitializeAccountRequest) (*model.InitializeAccountResult, error) {
	var initializeAccountResult []*model.InitializeAccountResult

	err := pgxscan.Select(context.Background(), handler.Pool, &initializeAccountResult, initializeAccountSql, request.PortfolioId, request.UserId, request.Currency)
	if err != nil {
		return nil, handler.handleErrors(err)
	}

	return initializeAccountResult[0], nil
}

func (handler *PostgresRepository) CreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (*model.CreateTransactionResult, error) {
	var createTransactionResult []*model.CreateTransactionResult

	transactionType, ok := utils.GetStringFromTransactionType(request.TransactionType)
	if !ok {
		return nil, fmt.Errorf("bad request: transaction type not supported: %v", transactionType)
	}

	totalAmountInt, _ := strconv.ParseInt(request.TotalAmount, 10, 64)
	if request.FeeAmount != nil {
		feeAmountInt, _ := strconv.ParseInt(request.FeeAmount.Value, 10, 64)
		totalAmountInt += feeAmountInt
	}

	err := pgxscan.Select(context.Background(),
		handler.Pool,
		&createTransactionResult,
		createTransactionSql,
		request.OrderId,
		request.Sender.Currency,
		request.Sender.UserId,
		request.Receiver.Currency,
		request.Receiver.UserId,
		request.RequestId.Value,
		totalAmountInt,
		transactionType)
	if err != nil {
		return nil, handler.handleErrors(err)
	}

	return createTransactionResult[0], nil
}

func (handler *PostgresRepository) PartialReleaseHold(ctx context.Context, request *api.PartialReleaseHoldRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	var retailFeeAmount int64
	if request.RetailFeeAmount == nil {
		retailFeeAmount = 0
	} else {
		retailFeeAmount, _ = strconv.ParseInt(request.RetailFeeAmount.Value, 10, 64)
	}

	var venueFeeAmount int64
	if request.VenueFeeAmount == nil {
		venueFeeAmount = 0
	} else {
		venueFeeAmount, _ = strconv.ParseInt(request.VenueFeeAmount.Value, 10, 64)
	}

	retailAccountId, venueAccountId := utils.GetFeeAccounts("USD")

	err := pgxscan.Select(
		context.Background(),
		handler.Pool,
		&TransactionResult,
		partialReleaseHoldSql,
		request.OrderId,
		request.RequestId,
		request.SenderAmount,
		request.ReceiverAmount,
		retailFeeAmount,
		retailAccountId,
		venueFeeAmount,
		venueAccountId)
	if err != nil {
		return nil, handler.handleErrors(err)
	}

	return TransactionResult[0], nil
}

func (handler *PostgresRepository) CompleteTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	err := pgxscan.Select(
		context.Background(),
		handler.Pool,
		&TransactionResult,
		completeTransactionSql,
		request.OrderId,
		request.RequestId)
	if err != nil {
		return nil, handler.handleErrors(err)
	}

	return TransactionResult[0], nil
}

func (handler *PostgresRepository) FailTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	err := pgxscan.Select(context.Background(), handler.Pool, &TransactionResult, failTransactionSql, request.OrderId, request.RequestId)
	if err != nil {
		return nil, handler.handleErrors(err)
	}

	return TransactionResult[0], nil
}

func (handler *PostgresRepository) CancelTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	err := pgxscan.Select(context.Background(), handler.Pool, &TransactionResult, cancelTransactionSql, request.OrderId, request.RequestId)
	if err != nil {
		return nil, handler.handleErrors(err)
	}

	return TransactionResult[0], nil
}

func (handler *PostgresRepository) GetAllAccountsAndMostRecentBalances(ctx context.Context, userId string) ([]*model.GetAccountResult, error) {
	var accountResult []*model.GetAccountResult

	err := pgxscan.Select(context.Background(), handler.Pool, &accountResult, getAllAccountsAndMostRecentBalancesSql, userId)
	if err != nil {
		return nil, handler.handleErrors(err)
	}
	return accountResult, nil
}
