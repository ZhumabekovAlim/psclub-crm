package models

type PriceItem struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	CategoryID      int     `json:"category_id"`
	SubcategoryID   int     `json:"subcategory_id"`
	Quantity        int     `json:"quantity"`
	SalePrice       float64 `json:"sale_price"`
	BuyPrice        float64 `json:"buy_price"`
	IsSet           bool    `json:"is_set"`                     // true если это сет (комплект товаров)
	SubcategoryName string  `json:"subcategory_name,omitempty"` // имя подкатегории, если есть
}
