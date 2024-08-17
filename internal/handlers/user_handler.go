package handlers

import (
	"context"
	"github.com/moonicy/gofermart/internal/models"
)

type UsersStorage interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, login string) (models.User, error)
	SetToken(ctx context.Context, user models.User) error
}

type UsersHandler struct {
	usersStorage UsersStorage
}

func NewUsersHandler(usersStorage UsersStorage) *UsersHandler {
	return &UsersHandler{usersStorage}
}
