package handlers

import (
	"context"
	"github.com/moonicy/gofermart/internal/models"
)

//go:generate mockery --output ../../mocks --filename movement_storage_mock_gen.go --outpkg mocks --name MovementsStorage --with-expecter
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
