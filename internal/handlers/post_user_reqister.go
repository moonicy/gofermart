package handlers

import (
	"encoding/json"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/pkg/hash"
	"io"
	"net/http"
	"time"
)

func (uh *UsersHandler) PostUserReqister(res http.ResponseWriter, req *http.Request) {
	var user models.User

	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	if err = json.Unmarshal(body, &user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	hashpass, err := hash.HashPassword(user.Password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	user.Password = hashpass
	user.AuthToken = hash.MakeToken(user.Login)
	user.AuthTokenExpired = time.Now().Add(time.Hour * 24)

	err = uh.usersStorage.CreateUser(req.Context(), user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "Authorization",
		Value:    user.AuthToken,
		Expires:  user.AuthTokenExpired,
		Secure:   true,
		HttpOnly: true,
	})

}
