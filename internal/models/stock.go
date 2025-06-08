package models

import "time"

type StockHistory struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	ItemID      int       `json:"item_id"`
	Quantity    int       `json:"quantity"`
	BuyPrice    float64   `json:"buy_price"`
	TotalPrice  float64   `json:"total_price"`
	UserID      int       `json:"user_id"`
	Description string    `json:"description,omitempty"`
}
