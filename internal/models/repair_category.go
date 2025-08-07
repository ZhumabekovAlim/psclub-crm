package models

type RepairCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CompanyID int    `json:"company_id"`
	BranchID  int    `json:"branch_id"`
}
