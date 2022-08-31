package repository

import (
	"LedgerApp/model"
	"LedgerApp/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"

	api "LedgerApp/protos/ledger"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func NewPostgresHandler() *PostgresRepository {
	dbCredsString := os.Getenv("DB_CREDENTIALS")
	if dbCredsString == "" {
		log.Fatalf("no environment variable set for DB_CREDENTIALS")
	}

	var dbCredsJson map[string]interface{}
	json.Unmarshal([]byte(dbCredsString), &dbCredsJson)

	dbUsername := dbCredsJson["username"].(string)
	dbPassword := dbCredsJson["password"].(string)

	dbEndpoint := os.Getenv("DB_HOSTNAME")
	if dbEndpoint == "" {
		log.Fatalf("no environment variable set for DB_HOSTNAME")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatalf("no environment variable set for DB_PORT")
	}

	pool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/ledger", dbUsername, dbPassword, dbEndpoint, dbPort))

	if err != nil {
		log.Errorf("Failed to establish DB Pool: %v", err)
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

	const sql = `SELECT id, sender_id, receiver_id, request_id, transaction_type, created_at FROM create_transaction_and_place_hold($1, $2, $3, $4, $5, $6)`

	err := pgxscan.Select(context.Background(), handler.Pool, &createTransactionResult, sql, request.OrderId, request.SenderId, request.ReceiverId, request.RequestId, request.SenderAmount, utils.GetStringFromTransactionType(request.TransactionType))
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
