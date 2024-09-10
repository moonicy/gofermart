package handlers

import (
	"context"
	"github.com/moonicy/gofermart/internal/models"
)

//go:generate mockery --output ../../mocks --filename user_storage_mock_gen.go --outpkg mocks --name UsersStorage --with-expecter
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
