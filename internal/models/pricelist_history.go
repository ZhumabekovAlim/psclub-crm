package models

import "time"

type PricelistHistory struct {
	ID          int       `json:"id"`
	PriceItemID int       `json:"price_item_id"`
	Quantity    int       `json:"quantity"`
	BuyPrice    float64   `json:"buy_price"`
	Total       float64   `json:"total"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}
