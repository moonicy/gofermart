package handlers

import "github.com/moonicy/gofermart/internal/storage"

type MovementHandler struct {
	movementsStorage *storage.MovementsStorage
}

func NewMovementHandler(movementsStorage *storage.MovementsStorage) *MovementHandler {
	return &MovementHandler{movementsStorage}
}
