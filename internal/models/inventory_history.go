package models

import "time"

type InventoryHistory struct {
	ID          int       `json:"id"`
	PriceItemID int       `json:"price_item_id"`
	Expected    float64   `json:"expected"`
	Actual      float64   `json:"actual"`
	Difference  float64   `json:"difference"`
	CreatedAt   time.Time `json:"created_at"`
}
