package handlers

import (
	"encoding/json"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/pkg/hash"
	"io"
	"net/http"
	"time"
)

func (us *UsersHandler) PostUserLogin(res http.ResponseWriter, req *http.Request) {
	var user models.User

	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	if err = json.Unmarshal(body, &user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	foundUser, err := us.usersStorage.GetUser(req.Context(), user.Login)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	if !hash.VerifyPassword(user.Password, foundUser.Password) {
		http.Error(res, err.Error(), http.StatusUnauthorized)
	}
	token := hash.MakeToken(user.Login)
	expiredAt := time.Now().Add(time.Hour * 24)
	user.AuthToken = token
	user.AuthTokenExpired = expiredAt

	err = us.usersStorage.SetToken(req.Context(), user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  expiredAt,
		Secure:   true,
		HttpOnly: true,
	})
}
