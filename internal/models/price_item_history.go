package models

import "time"

type PriceListAddItem struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	ItemName  string    `json:"item_name"`
	Quantity  int       `json:"quantity"`
	BuyPrice  float64   `json:"buy_price"`
	Total     float64   `json:"total"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
