package models

type Category struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CompanyID int    `json:"company_id"`
	BranchID  int    `json:"branch_id"`
}
type Subcategory struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
	CompanyID  int    `json:"company_id"`
	BranchID   int    `json:"branch_id"`
}
