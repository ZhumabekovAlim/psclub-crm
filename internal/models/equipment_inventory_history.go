package models

import "time"

type EquipmentInventoryHistory struct {
	ID          int       `json:"id"`
	EquipmentID int       `json:"equipment_id"`
	Name        string    `json:"name,omitempty"`
	Expected    float64   `json:"expected"`
	Actual      float64   `json:"actual"`
	Difference  float64   `json:"difference"`
	CreatedAt   time.Time `json:"created_at"`
	CompanyID   int       `json:"company_id"`
	BranchID    int       `json:"branch_id"`
}
