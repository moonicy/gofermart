package models

import "time"

type Movement struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Number      string    `json:"number"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
