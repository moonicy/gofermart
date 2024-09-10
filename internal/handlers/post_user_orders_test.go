package handlers

import (
	"context"
	"github.com/moonicy/gofermart/internal/contextkey"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/internal/storage"
	"github.com/moonicy/gofermart/mocks"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOrdersHandler_PostUserOrders(t *testing.T) {
	tests := []struct {
		name        string
		number      string
		userID      int
		foundUserID int
		err         error
		code        int
	}{
		{name: "Create new order", number: "12345678903", userID: 123, err: storage.ErrNotFound, code: http.StatusAccepted},
		{name: "Order already exist on this user", number: "12345678903", userID: 123, foundUserID: 123, code: http.StatusOK},
		{name: "Order already exist on another user", number: "12345678903", userID: 123, foundUserID: 456, code: http.StatusConflict},
		{name: "Incorrect order number", number: "1", userID: 123, code: http.StatusUnprocessableEntity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := mocks.NewOrdersStorage(t)
			os.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(nil).Maybe()

			var foundOrder models.Order
			if tt.foundUserID != 0 {
				foundOrder = models.Order{UserID: tt.foundUserID}
			}
			os.EXPECT().GetOrder(mock.Anything, tt.number).Return(foundOrder, tt.err).Maybe()
			uh := &OrdersHandler{
				ordersStorage: os,
			}

			res := httptest.NewRecorder()
			req := &http.Request{}
			req = req.WithContext(context.WithValue(req.Context(), contextkey.UserKey, models.User{ID: tt.userID}))

			req.Body = io.NopCloser(strings.NewReader(tt.number))

			uh.PostUserOrders(res, req)

			if tt.code != res.Code {
				t.Errorf("got %d, want %d", res.Code, tt.code)
			}
		})
	}
}
