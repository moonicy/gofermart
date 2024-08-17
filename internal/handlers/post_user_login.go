package handlers

import (
	"encoding/json"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/pkg/hash"
	"io"
	"net/http"
	"time"
)

func (uh *UsersHandler) PostUserLogin(res http.ResponseWriter, req *http.Request) {
	var ur UserRequest

	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	if err = json.Unmarshal(body, &ur); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	err = ur.Validate()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	foundUser, err := uh.usersStorage.GetUser(req.Context(), ur.Login)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	if !hash.VerifyPassword(ur.Password, foundUser.Password) {
		http.Error(res, "Password incorrect", http.StatusUnauthorized)
	}

	user := models.User{
		Login:            ur.Login,
		Password:         ur.Password,
		AuthToken:        hash.MakeToken(ur.Login),
		AuthTokenExpired: time.Now().Add(time.Hour * 24),
	}

	err = uh.usersStorage.SetToken(req.Context(), user)
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
