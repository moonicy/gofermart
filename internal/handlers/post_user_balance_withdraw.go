package handlers

import (
	"encoding/json"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/moonicy/gofermart/internal/models"
	"io"
	"net/http"
)

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (mh *MovementHandler) PostUserBalanceWithdraw(res http.ResponseWriter, req *http.Request) {
	var movement models.Movement
	var wr WithdrawRequest

	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &wr)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if wr.Sum < 0 {
		http.Error(res, "negative amount", http.StatusUnprocessableEntity)
		return
	}

	if err = goluhn.Validate(wr.Order); err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	movement.Number = wr.Order
	movement.Sum = wr.Sum

	user := models.GetUserFromContext(req.Context())
	movement.UserID = user.ID

	if user.Accrual < movement.Sum {
		http.Error(res, "account not enough", http.StatusPaymentRequired)
		return
	}

	err = mh.movementsStorage.MakeWithdraw(req.Context(), movement)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
