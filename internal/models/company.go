package models

type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Branch struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"company_id"`
	Name      string `json:"name"`
}
