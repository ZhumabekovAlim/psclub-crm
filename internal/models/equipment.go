package models

type Equipment struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    float64 `json:"quantity"`
	Description string  `json:"description,omitempty"`
	CompanyID   int     `json:"company_id"`
	BranchID    int     `json:"branch_id"`
}
