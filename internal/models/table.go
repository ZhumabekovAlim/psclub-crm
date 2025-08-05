package models

type Table struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
	Number     int    `json:"number"`
	CompanyID  int    `json:"company_id"`
	BranchID   int    `json:"branch_id"`
}
