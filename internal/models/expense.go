package models

import "time"

type Expense struct {
	ID               int       `json:"id"`
	Date             time.Time `json:"date"`
	Title            string    `json:"title"`
	CategoryID       int       `json:"category_id,omitempty"`
	Category         string    `json:"category"`
	RepairCategoryID int       `json:"repair_category_id,omitempty"`
	RepairCategory   string    `json:"repair_category"`
	Total            float64   `json:"total"`
	Description      string    `json:"description"`
	Paid             bool      `json:"paid"`
	CreatedAt        time.Time `json:"created_at"`
}
