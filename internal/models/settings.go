package models

type Settings struct {
	ID               int           `json:"id"`
	PaymentType      int           `json:"payment_type"`
	BlockTime        int           `json:"block_time"`
	BonusPercent     int           `json:"bonus_percent"`
	WorkTimeFrom     string        `json:"work_time_from"`
	WorkTimeTo       string        `json:"work_time_to"`
	TablesCount      int           `json:"tables_count"`
	NotificationTime int           `json:"notification_time"`
	PaymentTypes     []PaymentType `json:"payment_types"` // список всех типов
	Channels         []Channel     `json:"channels"`      // список всех каналов
}
