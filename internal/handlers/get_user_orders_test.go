package handlers

import (
	"context"
	"encoding/json"
	"github.com/moonicy/gofermart/internal/contextkey"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOrdersHandler_GetUserOrders(t *testing.T) {
	tests := []struct {
		name   string
		orders []models.Order
		userID int
		code   int
	}{
		{
			name: "Get orders",
			orders: []models.Order{
				{
					ID:         1,
					UserID:     123,
					Accrual:    111,
					Number:     "9278923470",
					Status:     "PROCESSED",
					UploadedAt: time.Time{},
				},
				{
					ID:         2,
					UserID:     123,
					Accrual:    500,
					Number:     "12345678903",
					Status:     "PROCESSED",
					UploadedAt: time.Time{},
				},
			},
			userID: 123,
			code:   http.StatusOK,
		},
		{
			name:   "No data for response",
			orders: []models.Order{},
			userID: 123,
			code:   http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := mocks.NewOrdersStorage(t)

			os.EXPECT().GetOrders(mock.Anything, tt.userID).Return(tt.orders, nil).Maybe()
			uh := &OrdersHandler{
				ordersStorage: os,
			}

			res := httptest.NewRecorder()
			req := &http.Request{}
			req = req.WithContext(context.WithValue(req.Context(), contextkey.UserKey, models.User{ID: tt.userID}))

			uh.GetUserOrders(res, req)

			if tt.code != res.Code {
				t.Errorf("got %d, want %d", res.Code, tt.code)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			var o OrdersResponse
			if err = json.Unmarshal(body, &o); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			assert.Equal(t, len(tt.orders), len(o))
		})
	}
}
