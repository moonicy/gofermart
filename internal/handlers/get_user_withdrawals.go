package handlers

import (
	"encoding/json"
	"github.com/moonicy/gofermart/internal/models"
	"net/http"
	"time"
)

type WithdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawalsResponse []WithdrawalResponse

func NewWithdrawalsResponse(m []models.Movement) WithdrawalsResponse {
	result := make([]WithdrawalResponse, 0, len(m))
	for _, val := range m {
		result = append(result, WithdrawalResponse{
			Order:       val.Number,
			Sum:         val.Sum,
			ProcessedAt: val.ProcessedAt,
		})
	}
	return result
}

func (mh *MovementHandler) GetUserWithdrawals(res http.ResponseWriter, req *http.Request) {
	user := models.GetUserFromContext(req.Context())

	res.Header().Set("Content-Type", "application/json")

	movements, err := mh.movementsStorage.GetMovements(req.Context(), user.ID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(movements) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	wr := NewWithdrawalsResponse(movements)

	out, err := json.Marshal(wr)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Write(out)
}
