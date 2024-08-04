package storage

import (
	"context"
	"database/sql"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (id serial PRIMARY KEY, login text UNIQUE NOT NULL, password text NOT NULL, accrual bigint default 0 NOT NULL, auth_token text, auth_token_expired timestamp);`)
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS orders (id serial PRIMARY KEY, number text UNIQUE NOT NULL, user_id int NOT NULL REFERENCES users(id), status text default 'REGISTERED' NOT NULL, accrual bigint default 0 NOT NULL, uploaded_at timestamp default now());`)
	if err != nil {
		return err
	}
	return nil
}
