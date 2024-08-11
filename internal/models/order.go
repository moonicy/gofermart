package models

import "time"

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Accrual    float64   `json:"accrual"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}
