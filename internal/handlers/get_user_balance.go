package handlers

import (
	"encoding/json"
	"github.com/moonicy/gofermart/internal/models"
	"log"
	"net/http"
)

type BalanceResponse struct {
	Accrual   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

func (uh *UsersHandler) GetUserBalance(res http.ResponseWriter, req *http.Request) {
	user := models.GetUserFromContext(req.Context())

	accrual, withdrawn, err := uh.usersStorage.GetBalance(req.Context(), user.AuthToken)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody := BalanceResponse{
		Accrual:   accrual,
		Withdrawn: withdrawn,
	}

	out, err := json.Marshal(resBody)
	if err != nil {
		log.Fatal(err)
	}
	res.Write(out)
}
