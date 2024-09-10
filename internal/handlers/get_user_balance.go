package handlers

import (
	"encoding/json"
	"github.com/moonicy/gofermart/internal/models"
	"log"
	"net/http"
)

type balanceResponse struct {
	Accrual   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (uh *UsersHandler) GetUserBalance(res http.ResponseWriter, req *http.Request) {
	user := models.GetUserFromContext(req.Context())

	res.Header().Set("Content-Type", "application/json")

	resBody := balanceResponse{
		Accrual:   user.Accrual,
		Withdrawn: user.Withdrawn,
	}

	out, err := json.Marshal(resBody)
	if err != nil {
		log.Fatal(err)
	}
	res.Write(out)
}
