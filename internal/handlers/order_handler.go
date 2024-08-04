package handlers

import "github.com/moonicy/gofermart/internal/storage"

type OrdersHandler struct {
	ordersStorage *storage.OrdersStorage
}

func NewOrdersHandler(ordersStorage *storage.OrdersStorage) *OrdersHandler {
	return &OrdersHandler{ordersStorage}
}
