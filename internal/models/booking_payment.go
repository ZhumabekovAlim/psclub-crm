package models

// BookingPayment represents a part of a booking payment with specific method.
type BookingPayment struct {
	ID            int     `json:"id"`
	BookingID     int     `json:"booking_id"`
	CompanyID     int     `json:"company_id"`
	BranchID      int     `json:"branch_id"`
	PaymentTypeID int     `json:"payment_type_id"`
	Amount        int     `json:"amount"`
	PaymentType   *string `json:"payment_type,omitempty"`
}
