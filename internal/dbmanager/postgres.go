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
