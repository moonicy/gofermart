package storage

import (
	"database/sql"
	"errors"
	"github.com/moonicy/gofermart/internal/config"
)

var ErrConflict = errors.New("entity already exists")
var ErrNotFound = errors.New("entity not found")

func NewDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	return db, nil
}
