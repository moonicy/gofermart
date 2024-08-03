package users

import "time"

type User struct {
	ID               int       `json:"id"`
	Login            string    `json:"login"`
	Password         string    `json:"password"`
	Accrual          int       `json:"accrual"`
	AuthToken        string    `json:"auth_token"`
	AuthTokenExpired time.Time `json:"auth_token_expired"`
}
