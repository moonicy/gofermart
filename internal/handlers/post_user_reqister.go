package handlers

import (
	"encoding/json"
	"errors"
	"github.com/moonicy/gofermart/internal/models"
	"github.com/moonicy/gofermart/pkg/hash"
	"io"
	"net/http"
	"time"
)

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (ur UserRequest) Validate() error {
	if ur.Login == "" {
		return errors.New("login is required")
	}
	if ur.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (uh *UsersHandler) PostUserReqister(res http.ResponseWriter, req *http.Request) {
	var ur UserRequest

	res.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(body, &ur); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = ur.Validate()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	hashpass, err := hash.HashPassword(ur.Password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	user := models.User{
		Login:            ur.Login,
		Password:         hashpass,
		AuthToken:        hash.MakeToken(ur.Login),
		AuthTokenExpired: time.Now().Add(time.Hour * 24),
	}

	err = uh.usersStorage.CreateUser(req.Context(), user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	http.SetCookie(res, &http.Cookie{
		Name:     "Authorization",
		Value:    user.AuthToken,
		Expires:  user.AuthTokenExpired,
		Secure:   false,
		HttpOnly: false,
	})

}
