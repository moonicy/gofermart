package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moonicy/gofermart/internal/models"
)

type OrdersStorage struct {
	db *sql.DB
}

func NewOrdersStorage(db *sql.DB) *OrdersStorage {
	return &OrdersStorage{db: db}
}

func (os *OrdersStorage) CreateOrder(ctx context.Context, order models.Order) error {
	_, err := os.db.ExecContext(ctx, `INSERT INTO orders (number, user_id) VALUES ($1, $2)`, order.Number, order.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
		return err
	}
	return nil
}

func (os *OrdersStorage) GetOrder(ctx context.Context, number string) (models.Order, error) {
	var order models.Order
	row := os.db.QueryRowContext(ctx, `SELECT id, number, user_id, status, accrual, uploaded_at FROM orders WHERE number = $1`, number)

	err := row.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Order{}, ErrNotFound
		}
		return models.Order{}, err
	}
	return order, nil
}
