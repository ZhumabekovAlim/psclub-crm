package models

import "time"

type Booking struct {
	ID             int           `json:"id"`
	ClientID       int           `json:"client_id"`
	TableID        int           `json:"table_id"`
	UserID         int           `json:"user_id"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Note           string        `json:"note"`
	Discount       int           `json:"discount"`
	DiscountReason string        `json:"discount_reason"`
	TotalAmount    int           `json:"total_amount"`
	BonusUsed      int           `json:"bonus_used"`
	PaymentStatus  string        `json:"payment_status"`
	PaymentTypeID  int           `json:"payment_type_id"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	ClientName     string        `json:"client_name,omitempty"`
	ClientPhone    string        `json:"client_phone,omitempty"`
	PaymentType    *string       `json:"payment_type,omitempty"`
	Items          []BookingItem `json:"items,omitempty"`
}

type BookingItem struct {
	ID        int     `json:"id"`
	BookingID int     `json:"booking_id"`
	ItemID    int     `json:"item_id"`
	Quantity  int     `json:"quantity"`
	Price     int     `json:"price"`
	Discount  int     `json:"discount"`
	ItemPrice float64 `json:"item_price,omitempty"`
}
