package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Password     string    `json:"password"`
	Role         string    `json:"role"`
	Permissions  []string  `json:"permissions"`
	SalaryHookah float64   `json:"salary_hookah"`
	SalaryBar    float64   `json:"salary_bar"`
	SalaryShift  int       `json:"salary_shift"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
