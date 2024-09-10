package handlers

import (
	"context"
	"fmt"
	"github.com/moonicy/gofermart/internal/contextkey"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/mocks"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMovementHandler_PostUserBalanceWithdraw(t *testing.T) {
	tests := []struct {
		name     string
		movement models.Movement
		user     models.User
		code     int
	}{
		{name: "Get movements", movement: models.Movement{Number: "9278923470", Sum: 123.1}, user: models.User{ID: 123, Accrual: 222}, code: http.StatusOK},
		{name: "Not enough for payment", movement: models.Movement{Number: "12345678903", Sum: 123.1}, user: models.User{ID: 123, Accrual: 5}, code: http.StatusPaymentRequired},
		{name: "Incorrect order number", movement: models.Movement{Number: "1", Sum: 123.1}, user: models.User{ID: 123, Accrual: 222}, code: http.StatusUnprocessableEntity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := mocks.NewMovementsStorage(t)

			ms.EXPECT().MakeWithdraw(mock.Anything, mock.Anything).Return(nil).Maybe()
			mh := &MovementHandler{
				movementsStorage: ms,
			}

			res := httptest.NewRecorder()
			req := &http.Request{}
			req = req.WithContext(context.WithValue(req.Context(), contextkey.UserKey, models.User{ID: tt.user.ID, Accrual: tt.user.Accrual}))

			body := fmt.Sprintf("{\"order\":\"%s\",\"sum\":%.1f}", tt.movement.Number, tt.movement.Sum)
			req.Body = io.NopCloser(strings.NewReader(body))

			mh.PostUserBalanceWithdraw(res, req)

			if tt.code != res.Code {
				t.Errorf("got %d, want %d", res.Code, tt.code)
			}
		})
	}
}
