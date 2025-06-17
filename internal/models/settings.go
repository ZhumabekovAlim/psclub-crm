package models

type Settings struct {
	ID           int           `json:"id"`
	PaymentType  int           `json:"payment_type"`
	BlockTime    int           `json:"block_time"`
	BonusPercent int           `json:"bonus_percent"`
	WorkTimeFrom string        `json:"work_time_from"`
	WorkTimeTo   string        `json:"work_time_to"`
	PaymentTypes []PaymentType `json:"payment_types"` // список всех типов
}
