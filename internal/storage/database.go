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

type Tx interface {
	Begin() (*sql.Tx, error)
}

type TxDB interface {
	DB
	Tx
}

func NewDB(cfg config.Config) (TxDB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDBorTX(ctx context.Context, db DB) DB {
	val := ctx.Value(contextkey.TransactionKey)
	tx, ok := val.(*sql.Tx)
	if ok {
		return tx
	} else {
		return db
	}
}
