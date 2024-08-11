package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/moonicy/gofermart/internal/contextkey"
	"github.com/moonicy/gofermart/internal/models"
)

type UserOrderStorage struct {
	db *sql.DB
	os *OrdersStorage
	us *UsersStorage
}

func NewUserOrderStorage(db *sql.DB, os *OrdersStorage, us *UsersStorage) *UserOrderStorage {
	return &UserOrderStorage{db: db, os: os, us: us}
}

func (uos *UserOrderStorage) UpdateAccrual(ctx context.Context, order models.Order) error {
	tx, err := uos.db.Begin()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, contextkey.TransactionKey, tx)
	err = uos.os.UpdateOrder(ctx, order)
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return errors.Join(err, errRb)
		}
		return err
	}
	err = uos.us.AddAccrual(ctx, order.UserID, order.Accrual)
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			return errors.Join(err, errRb)
		}
		return err
	}
	return tx.Commit()
}
