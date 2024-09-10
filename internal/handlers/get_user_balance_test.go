package handlers

import (
	"context"
	"fmt"
	"github.com/moonicy/gofermart/internal/contextkey"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersHandler_GetUserBalance(t *testing.T) {
	tests := []struct {
		name string
		user models.User
		err  error
		code int
	}{
		{name: "Get balance", user: models.User{Accrual: 123.2, Withdrawn: 123}, code: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := mocks.NewUsersStorage(t)

			uh := &UsersHandler{
				usersStorage: us,
			}

			res := httptest.NewRecorder()
			req := &http.Request{}
			req = req.WithContext(context.WithValue(req.Context(), contextkey.UserKey, tt.user))

			body := fmt.Sprintf("{\"current\":%.1f,\"withdrawn\":%.0f}", tt.user.Accrual, tt.user.Withdrawn)

			uh.GetUserBalance(res, req)

			if tt.code != res.Code {
				t.Errorf("got %d, want %d", res.Code, tt.code)
			}

			r := res.Result()

			got, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, body, string(got))
		})
	}
}
