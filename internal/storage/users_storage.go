package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moonicy/gofermart/internal/users"
)

var ErrConflict = errors.New("user already exists")
var ErrNotFound = errors.New("user not found")

type UsersStorage struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) *UsersStorage {
	return &UsersStorage{db: db}
}

func (us *UsersStorage) GetUser(ctx context.Context, login string) (users.User, error) {
	var user users.User
	row := us.db.QueryRowContext(ctx, `SELECT id, login, password, accrual, auth_token, auth_token_expired FROM users WHERE login = $1`, login)
	var authToken sql.NullString
	var authTokenExpired sql.NullTime

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Accrual, &authToken, &authTokenExpired)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users.User{}, ErrNotFound
		}
		return users.User{}, err
	}
	if authToken.Valid {
		user.AuthToken = authToken.String
	}
	if authTokenExpired.Valid {
		user.AuthTokenExpired = authTokenExpired.Time
	}
	return user, nil
}

func (us *UsersStorage) SetUser(ctx context.Context, user users.User) error {
	_, err := us.db.ExecContext(ctx, `INSERT INTO users (login, password, auth_token, auth_token_expired) VALUES ($1, $2, $3, $4)`, user.Login, user.Password, user.AuthToken, user.AuthTokenExpired)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
		return err
	}
	return nil
}

func (us *UsersStorage) SetToken(ctx context.Context, user users.User) error {
	_, err := us.db.ExecContext(ctx, `UPDATE users SET auth_token = $1, auth_token_expired = $2 WHERE login = $3`, user.AuthToken, user.AuthTokenExpired, user.Login)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
		return err
	}
	return nil
}
