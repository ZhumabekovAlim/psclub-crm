package models

type PaymentType struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	HoldPercent float64 `json:"hold_percent"`
	CompanyID   int     `json:"company_id"`
	BranchID    int     `json:"branch_id"`
}
