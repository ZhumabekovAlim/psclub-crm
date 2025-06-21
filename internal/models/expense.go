package models

import (
	"database/sql"
	"time"
)

type Expense struct {
	ID          int           `json:"id"`
	Date        time.Time     `json:"date"`
	Title       string        `json:"title"`
	CategoryID  sql.NullInt64 `json:"category_id"`
	Category    string        `json:"category"`
	Total       float64       `json:"total"`
	Description string        `json:"description"`
	Paid        bool          `json:"paid"`
	CreatedAt   time.Time     `json:"created_at"`
}
