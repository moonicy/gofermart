package handlers

import (
	"encoding/json"
	"github.com/moonicy/gofermart/internal/models"
	"log"
	"net/http"
	"time"
)

type OrderResponse struct {
	Accrual    float64   `json:"accrual"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type OrdersResponse []OrderResponse

func NewOrdersResponse(o []models.Order) OrdersResponse {
	result := make([]OrderResponse, 0, len(o))
	for _, val := range o {
		result = append(result, OrderResponse{
			Accrual:    val.Accrual,
			Number:     val.Number,
			Status:     val.Status,
			UploadedAt: val.UploadedAt,
		})
	}
	return result
}

func (oh *OrdersHandler) GetUserOrders(res http.ResponseWriter, req *http.Request) {
	var orders []models.Order

	user := models.GetUserFromContext(req.Context())

	orders, err := oh.ordersStorage.GetOrders(req.Context(), user.ID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if orders == nil {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	result := NewOrdersResponse(orders)

	out, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	res.Write(out)
}
