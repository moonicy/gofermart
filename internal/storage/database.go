package storage

import (
	"database/sql"
	"github.com/moonicy/gofermart/internal/config"
)

func NewDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	return db, nil
}
