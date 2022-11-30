package dbmanager

import "context"

type DBManager interface {
    Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

