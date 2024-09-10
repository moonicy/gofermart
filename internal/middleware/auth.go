package middleware

import (
	"context"
	"errors"
	"github.com/moonicy/gofermart/internal/contextkey"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/internal/storage"
	"net/http"
	"time"
)

type UsersStorage interface {
	GetUserByAuth(ctx context.Context, token string) (models.User, error)
}

func Auth(us UsersStorage) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			cookie, err := req.Cookie("Authorization")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			user, err := us.GetUserByAuth(req.Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, storage.ErrNotFound) {
					res.WriteHeader(http.StatusUnauthorized)
					return
				}
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			if user.AuthTokenExpired.Before(time.Now()) {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(req.Context(), contextkey.UserKey, user)
			handler.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
