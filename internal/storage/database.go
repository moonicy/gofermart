package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/moonicy/gofermart/internal/config"
	"github.com/moonicy/gofermart/internal/contextkey"
)

var ErrConflict = errors.New("entity already exists")
var ErrNotFound = errors.New("entity not found")

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDBorTX(ctx context.Context, db *sql.DB) DB {
	val := ctx.Value(contextkey.TransactionKey)
	tx, ok := val.(*sql.Tx)
	if ok {
		return tx
	} else {
		return db
	}
}
