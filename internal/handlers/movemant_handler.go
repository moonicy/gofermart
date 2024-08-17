package handlers

import (
	"context"
	"github.com/moonicy/gofermart/internal/models"
)

type MovementsStorage interface {
	GetMovements(ctx context.Context, userID int) ([]models.Movement, error)
	MakeWithdraw(ctx context.Context, movement models.Movement) error
}

type MovementHandler struct {
	movementsStorage MovementsStorage
}

func NewMovementHandler(movementsStorage MovementsStorage) *MovementHandler {
	return &MovementHandler{movementsStorage}
}
