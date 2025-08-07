package models

type Channel struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"company_id"`
	BranchID  int    `json:"branch_id"`
	Name      string `json:"name"`
}
