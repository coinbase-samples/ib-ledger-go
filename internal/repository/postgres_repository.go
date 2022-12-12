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
	"errors"
	"fmt"
	"math/big"

	"github.com/coinbase-samples/ib-ledger-go/internal/dbmanager"
	ledgererr "github.com/coinbase-samples/ib-ledger-go/internal/errors"
	"github.com/coinbase-samples/ib-ledger-go/internal/model"
	"github.com/coinbase-samples/ib-ledger-go/internal/utils"
	"google.golang.org/grpc/codes"

	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	initializeAccountSql = `
    SELECT 
        id, 
        portfolio_id, 
        user_id, 
        currency, 
        created_at, 
        balance, 
        hold,
        available 
    FROM initialize_account($1, $2, $3);`

	getAllAccountsAndMostRecentBalancesSql = `
    SELECT 
        account_id, 
        currency, 
        balance, 
        hold, 
        available, 
        created_at 
    FROM get_balances_for_users($1);`

	createTransactionSql = `
    SELECT 
        id, 
        sender_id, 
        receiver_id, 
        request_id, 
        transaction_type, 
        created_at 
    FROM create_transaction_and_place_hold($1, $2, $3, $4, $5, $6, $7, $8);`

	partialReleaseHoldSql = `
    SELECT 
        hold_id, 
        sender_entry_id, 
        receiver_entry_id, 
        sender_balance_id, 
        receiver_balance_id 
    FROM partial_release_hold($1, $2, $3, $4, $5, $6, $7, $8);`

	completeTransactionSql = `
    SELECT 
        hold_id, 
        sender_entry_id, 
        receiver_entry_id, 
        sender_balance_id, 
        receiver_balance_id 
    FROM complete_transaction($1, $2);`

	failTransactionSql = `
    SELECT 
        hold_id, 
        sender_entry_id, 
        receiver_entry_id, 
        sender_balance_id, 
        receiver_balance_id 
    FROM fail_transaction($1, $2);`

	cancelTransactionSql = `
    SELECT 
        hold_id, 
        sender_entry_id, 
        receiver_entry_id, 
        sender_balance_id, 
        receiver_balance_id 
    FROM cancel_transaction($1, $2);`
)

type PostgresRepository struct {
	DBManager dbmanager.DBManager
}

func NewPostgresHandler(dbmanager dbmanager.DBManager) *PostgresRepository {
	return &PostgresRepository{DBManager: dbmanager}
}

func (r *PostgresRepository) handleErrors(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		msg := pgErr.Message
		switch msg {
		case "LGR401":
			return ledgererr.New(codes.NotFound, "balance not found")
		case "LGR402":
			return ledgererr.New(codes.NotFound, "sender account not found")
		case "LGR403":
			return ledgererr.New(codes.NotFound, "receiver account not found")
		case "LGR404":
			return ledgererr.New(codes.NotFound, msg)
		case "LGR405":
			return ledgererr.New(codes.InvalidArgument, "no unreleased hold found for transaction")
		case "LGR501":
			return ledgererr.New(codes.InvalidArgument, "insufficient balance for transaction")
		case "LGR502":
			return ledgererr.New(codes.FailedPrecondition, "all releases have not been ledgered")
		default:
			return ledgererr.New(codes.Internal, "postgres internal failure")
		}
	}
	return ledgererr.New(codes.Internal, err.Error())
}

func (pr *PostgresRepository) InitializeAccount(ctx context.Context, request *api.InitializeAccountRequest) (*model.InitializeAccountResult, error) {
	var initializeAccountResult []*model.InitializeAccountResult

	err := pr.DBManager.Query(
		context.Background(),
		&initializeAccountResult,
		initializeAccountSql,
		request.PortfolioId,
		request.UserId,
		request.Currency)
	if err != nil {
		return nil, pr.handleErrors(err)
	}

	return initializeAccountResult[0], nil
}

func (pr *PostgresRepository) CreateTransaction(ctx context.Context, request *api.CreateTransactionRequest) (*model.CreateTransactionResult, error) {
	var createTransactionResult []*model.CreateTransactionResult

	transactionType, ok := utils.GetStringFromTransactionType(request.TransactionType)
	if !ok {
		return nil, fmt.Errorf("bad request: transaction type not supported: %v", transactionType)
	}

	totalAmountInt := new(big.Int)
	totalAmountInt.SetString(request.TotalAmount, 10)
	totalAmountWithFees := new(big.Int)
	if request.FeeAmount != nil {
		feeAmountInt := new(big.Int)
		feeAmountInt.SetString(request.FeeAmount.Value, 10)
		totalAmountWithFees.Add(totalAmountInt, feeAmountInt)
	} else {
		totalAmountWithFees = totalAmountInt
	}

	err := pr.DBManager.Query(context.Background(),
		&createTransactionResult,
		createTransactionSql,
		request.OrderId,
		request.Sender.Currency,
		request.Sender.UserId,
		request.Receiver.Currency,
		request.Receiver.UserId,
		request.RequestId.Value,
		totalAmountWithFees.String(),
		transactionType)
	if err != nil {
		return nil, pr.handleErrors(err)
	}

	return createTransactionResult[0], nil
}

func (pr *PostgresRepository) PartialReleaseHold(ctx context.Context, request *api.PartialReleaseHoldRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	var retailFeeAmount string
	if request.RetailFeeAmount == nil {
		retailFeeAmount = "0"
	} else {
		retailFeeAmount = request.RetailFeeAmount.Value
	}

	var venueFeeAmount string
	if request.VenueFeeAmount == nil {
		venueFeeAmount = "0"
	} else {
		venueFeeAmount = request.VenueFeeAmount.Value
	}

	retailAccountId, venueAccountId, err := utils.GetFeeAccounts("USD")
	if err != nil {
		return nil, err
	}

	err = pr.DBManager.Query(
		context.Background(),
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
		return nil, pr.handleErrors(err)
	}

	return TransactionResult[0], nil
}

func (pr *PostgresRepository) CompleteTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	err := pr.DBManager.Query(
		context.Background(),
		&TransactionResult,
		completeTransactionSql,
		request.OrderId,
		request.RequestId)
	if err != nil {
		return nil, pr.handleErrors(err)
	}

	return TransactionResult[0], nil
}

func (pr *PostgresRepository) FailTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	err := pr.DBManager.Query(
		context.Background(),
		&TransactionResult,
		failTransactionSql,
		request.OrderId,
		request.RequestId)
	if err != nil {
		return nil, pr.handleErrors(err)
	}

	return TransactionResult[0], nil
}

func (pr *PostgresRepository) CancelTransaction(ctx context.Context, request *api.FinalizeTransactionRequest) (*model.TransactionResult, error) {
	var TransactionResult []*model.TransactionResult

	err := pr.DBManager.Query(
		context.Background(),
		&TransactionResult,
		cancelTransactionSql,
		request.OrderId,
		request.RequestId)
	if err != nil {
		return nil, pr.handleErrors(err)
	}

	return TransactionResult[0], nil
}

func (pr *PostgresRepository) GetAllAccountsAndMostRecentBalances(ctx context.Context, userId string) ([]*model.GetAccountResult, error) {
	var accountResult []*model.GetAccountResult

	err := pr.DBManager.Query(
		context.Background(),
		&accountResult,
		getAllAccountsAndMostRecentBalancesSql,
		userId)
	if err != nil {
		return nil, pr.handleErrors(err)
	}
	return accountResult, nil
}
