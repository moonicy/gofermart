package handlers

import (
	"context"
	"github.com/moonicy/gofermart/internal/models"
)

//go:generate mockery --output ../../mocks --filename order_storage_mock_gen.go --outpkg mocks --name OrdersStorage --with-expecter
type OrdersStorage interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrder(ctx context.Context, number string) (models.Order, error)
	GetOrders(ctx context.Context, userID int) ([]models.Order, error)
}

type OrdersHandler struct {
	ordersStorage OrdersStorage
}

func NewOrdersHandler(ordersStorage OrdersStorage) *OrdersHandler {
	return &OrdersHandler{ordersStorage}
}
