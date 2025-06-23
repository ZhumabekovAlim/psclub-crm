package models

type PriceSet struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	CategoryID      int       `json:"category_id"`
	SubcategoryID   int       `json:"subcategory_id"`
	Price           int       `json:"price"`
	SubcategoryName string    `json:"subcategory_name,omitempty"`
	Items           []SetItem `json:"items"`
}

type SetItem struct {
	ID         int    `json:"id"`
	PriceSetID int    `json:"price_set_id"`
	ItemID     int    `json:"item_id"`
	Quantity   int    `json:"quantity"`
	ItemName   string `json:"item_name,omitempty"`
}
