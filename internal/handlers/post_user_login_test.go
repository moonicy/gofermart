package handlers

import (
	"fmt"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/mocks"
	"github.com/moonicy/gofermart/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUsersHandler_PostUserLogin(t *testing.T) {
	hashPass, err := hash.HashPassword("test")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		login      string
		password   string
		user       models.User
		getUserErr error
		setUserErr error
		code       int
	}{
		{name: "Login user", login: "test", password: "test", user: models.User{Login: "test", Password: hashPass, AuthToken: ""}, code: http.StatusOK},
		{name: "User already authorized", login: "test", password: "test", user: models.User{Login: "test", Password: hashPass, AuthToken: "123"}, code: http.StatusOK},
		{name: "Password incorrect", login: "test", password: "test", user: models.User{Login: "test", Password: "password", AuthToken: ""}, code: http.StatusUnauthorized},
		{name: "No login", login: "", password: "test", code: http.StatusBadRequest},
		{name: "No password", login: "test", password: "", code: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := mocks.NewUsersStorage(t)
			us.EXPECT().GetUser(mock.Anything, tt.login).Return(tt.user, tt.getUserErr).Maybe()
			us.EXPECT().SetToken(mock.Anything, mock.Anything).Return(tt.setUserErr).Maybe()
			uh := &UsersHandler{
				usersStorage: us,
			}

			res := httptest.NewRecorder()
			req := &http.Request{}

			body := fmt.Sprintf("{\"login\":\"%s\",\"password\":\"%s\"}", tt.login, tt.password)
			req.Body = io.NopCloser(strings.NewReader(body))

			uh.PostUserLogin(res, req)

			if tt.code != res.Code {
				t.Errorf("got %d, want %d", res.Code, tt.code)
			}
		})
	}
}
