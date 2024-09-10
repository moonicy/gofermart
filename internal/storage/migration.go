package storage

import (
	"context"
)

type Migrator struct {
	db DB
}

func NewMigrator(db DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate(ctx context.Context) error {
	_, err := m.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (id serial PRIMARY KEY, login text UNIQUE NOT NULL, password text NOT NULL, accrual double precision default 0 NOT NULL, withdrawn double precision default 0 NOT NULL, auth_token text, auth_token_expired timestamp);`)
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS orders (id serial PRIMARY KEY, number text UNIQUE NOT NULL, user_id int NOT NULL REFERENCES users(id), status text default 'NEW' NOT NULL, accrual double precision default 0 NOT NULL, uploaded_at timestamp default now());`)
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS movements (id serial PRIMARY KEY, user_id int NOT NULL REFERENCES users(id), number text UNIQUE NOT NULL, sum double precision NOT NULL, processed_at timestamp default now());`)
	if err != nil {
		return err
	}
	return nil
}
