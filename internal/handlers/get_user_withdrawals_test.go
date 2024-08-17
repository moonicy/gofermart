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

func TestMovementHandler_GetUserWithdrawals(t *testing.T) {
	tests := []struct {
		name      string
		movements []models.Movement
		userID    int
		code      int
	}{
		{
			name: "Get movements",
			movements: []models.Movement{
				{
					UserID:      123,
					Sum:         111,
					Number:      "9278923470",
					ProcessedAt: time.Time{},
				},
				{
					UserID:      123,
					Sum:         500,
					Number:      "12345678903",
					ProcessedAt: time.Time{},
				},
			},
			userID: 123,
			code:   http.StatusOK,
		},
		{
			name:      "No data for response",
			movements: []models.Movement{},
			userID:    123,
			code:      http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := mocks.NewMovementsStorage(t)

			ms.EXPECT().GetMovements(mock.Anything, tt.userID).Return(tt.movements, nil).Maybe()
			mh := &MovementHandler{
				movementsStorage: ms,
			}

			res := httptest.NewRecorder()
			req := &http.Request{}
			req = req.WithContext(context.WithValue(req.Context(), contextkey.UserKey, models.User{ID: tt.userID}))

			mh.GetUserWithdrawals(res, req)

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
			assert.Equal(t, len(tt.movements), len(o))
		})
	}
}
