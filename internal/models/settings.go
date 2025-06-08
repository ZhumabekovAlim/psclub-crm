package models

type Settings struct {
	ID           int    `json:"id"`
	PaymentType  string `json:"payment_type"`
	BlockTime    int    `json:"block_time"` // мин. до редактирования/удаления брони
	BonusPercent int    `json:"bonus_percent"`
	WorkTimeFrom string `json:"work_time_from"` // "10:00"
	WorkTimeTo   string `json:"work_time_to"`   // "02:00"
}
