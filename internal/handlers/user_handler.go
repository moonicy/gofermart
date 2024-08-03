package handlers

import "github.com/moonicy/gofermart/internal/storage"

type UsersHandler struct {
	usersStorage *storage.UsersStorage
}

func NewUsersHandler(usersStorage *storage.UsersStorage) *UsersHandler {
	return &UsersHandler{usersStorage}
}
