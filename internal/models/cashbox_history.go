package models

import "time"

type CashboxHistory struct {
	ID        int       `json:"id"`
	Operation string    `json:"operation"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
