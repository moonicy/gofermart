package models

import (
	"context"
	"github.com/moonicy/gofermart/internal/contextkey"
	"time"
)

type User struct {
	ID               int       `json:"id"`
	Login            string    `json:"login"`
	Password         string    `json:"password"`
	Accrual          int       `json:"accrual"`
	AuthToken        string    `json:"auth_token"`
	AuthTokenExpired time.Time `json:"auth_token_expired"`
}

func GetUserFromContext(ctx context.Context) User {
	value := ctx.Value(contextkey.UserKey)
	user, ok := value.(User)
	if !ok {
		return User{}
	}
	return user
}
