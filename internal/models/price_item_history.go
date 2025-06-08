package models

import "time"

type PriceItemHistory struct {
	ID          int       `json:"id"`
	PriceItemID int       `json:"price_item_id"`
	Operation   string    `json:"operation"` // "INCOME", "OUTCOME", "ADJUST"
	Quantity    int       `json:"quantity"`
	BuyPrice    float64   `json:"buy_price"` // цена закупки (если есть)
	Total       float64   `json:"total"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}
