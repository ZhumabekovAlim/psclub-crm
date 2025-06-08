package models

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Subcategory struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
}
