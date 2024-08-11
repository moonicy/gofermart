package handlers

import (
	"errors"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/internal/storage"
	"io"
	"net/http"
)

func (oh *OrdersHandler) PostUserOrders(res http.ResponseWriter, req *http.Request) {
	var order models.Order

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	orderNum := string(body)

	err = goluhn.Validate(orderNum)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	order.Number = orderNum

	user := models.GetUserFromContext(req.Context())
	order.UserID = user.ID

	foundOrder, err := oh.ordersStorage.GetOrder(req.Context(), order.Number)
	if err == nil {
		if foundOrder.UserID == order.UserID {
			res.WriteHeader(http.StatusOK)
		} else {
			res.WriteHeader(http.StatusConflict)
		}
		return
	}

	if !errors.Is(err, storage.ErrNotFound) {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	err = oh.ordersStorage.CreateOrder(req.Context(), order)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}
