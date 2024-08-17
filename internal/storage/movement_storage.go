package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moonicy/gofermart/internal/contextkey"
	"github.com/moonicy/gofermart/internal/models"
)

type MovementsStorage struct {
	db *sql.DB
	us *UsersStorage
	os *OrdersStorage
}

func NewMovementsStorage(db *sql.DB, us *UsersStorage, os *OrdersStorage) *MovementsStorage {
	return &MovementsStorage{db: db, us: us, os: os}
}

func (ms *MovementsStorage) CreateMovement(ctx context.Context, movement models.Movement) error {
	_, err := ms.db.ExecContext(ctx, `INSERT INTO movements (number, user_id, sum) VALUES ($1, $2, $3)`, movement.Number, movement.UserID, movement.Sum)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
		return err
	}
	return nil
}

func (ms *MovementsStorage) GetMovements(ctx context.Context, userID int) ([]models.Movement, error) {
	var movement models.Movement
	var movements []models.Movement
	row, err := ms.db.QueryContext(ctx, `SELECT id, number, user_id, sum, processed_at FROM movements WHERE user_id = $1`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return movements, ErrNotFound
		}
	}

	if row.Err() != nil {
		return movements, row.Err()
	}

	for row.Next() {
		err = row.Scan(&movement.ID, &movement.Number, &movement.UserID, &movement.Sum, &movement.ProcessedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []models.Movement{}, ErrNotFound
			}
			return []models.Movement{}, err
		}
		movements = append(movements, movement)
	}
	return movements, nil
}

func (ms *MovementsStorage) MakeWithdraw(ctx context.Context, movement models.Movement) error {
	tx, err := ms.db.Begin()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, contextkey.TransactionKey, tx)

	err = ms.os.CreateOrder(ctx, models.Order{
		UserID: movement.UserID,
		Number: movement.Number,
	})
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return errors.Join(err, errRb)
		}
		return err
	}

	err = ms.CreateMovement(ctx, movement)
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return errors.Join(err, errRb)
		}
		return err
	}

	err = ms.us.AddWithdraw(ctx, movement.UserID, movement.Sum)
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return errors.Join(err, errRb)
		}
		return err
	}
	return tx.Commit()
}
