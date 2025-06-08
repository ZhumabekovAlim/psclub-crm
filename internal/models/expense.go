package models

import "time"

type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Total       float64   `json:"total"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
