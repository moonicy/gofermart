package handlers

import (
	"errors"
	"fmt"
	"github.com/moonicy/gofermart/mocks"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUsersHandler_PostUserReqister(t *testing.T) {
	tests := []struct {
		name     string
		login    string
		password string
		err      error
		code     int
	}{
		{name: "Register new user", login: "test", password: "test", code: http.StatusOK},
		{name: "User already exist", login: "test", password: "test", err: errors.New("user already exist"), code: http.StatusInternalServerError},
		{name: "No login", login: "", password: "test", code: http.StatusBadRequest},
		{name: "No password", login: "test", password: "", code: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := mocks.NewUsersStorage(t)
			us.EXPECT().CreateUser(mock.Anything, mock.Anything).Return(tt.err).Maybe()
			uh := &UsersHandler{
				usersStorage: us,
			}

			res := httptest.NewRecorder()
			req := &http.Request{}

			body := fmt.Sprintf("{\"login\":\"%s\",\"password\":\"%s\"}", tt.login, tt.password)
			req.Body = io.NopCloser(strings.NewReader(body))

			uh.PostUserReqister(res, req)

			if tt.code != res.Code {
				t.Errorf("got %d, want %d", res.Code, tt.code)
			}
		})
	}
}
