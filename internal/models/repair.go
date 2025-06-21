package models

import "time"

type Repair struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	VIN         string    `json:"vin"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
