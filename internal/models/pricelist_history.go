package models

import "time"

type PricelistHistory struct {
	ID          int       `json:"id"`
	ItemName    string    `json:"item_name"`
	PriceItemID int       `json:"price_item_id"`
	Quantity    float64   `json:"quantity"`
	BuyPrice    float64   `json:"buy_price"`
	Total       float64   `json:"total"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UserName    string    `json:"user_name,omitempty"`
}
