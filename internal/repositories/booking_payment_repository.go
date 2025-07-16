package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

// BookingPaymentRepository handles CRUD for booking payments.
type BookingPaymentRepository struct {
	db *sql.DB
}

func NewBookingPaymentRepository(db *sql.DB) *BookingPaymentRepository {
	return &BookingPaymentRepository{db: db}
}

// Create inserts multiple booking payments for a booking.
func (r *BookingPaymentRepository) Create(ctx context.Context, bookingID int, payments []models.BookingPayment) error {
	if len(payments) == 0 {
		return nil
	}
	query := `INSERT INTO booking_payments (booking_id, payment_type_id, amount) VALUES (?, ?, ?)`
	for _, p := range payments {
		if _, err := r.db.ExecContext(ctx, query, bookingID, p.PaymentTypeID, p.Amount); err != nil {
			return err
		}
	}
	return nil
}

// DeleteByBookingID removes all payments for a booking.
func (r *BookingPaymentRepository) DeleteByBookingID(ctx context.Context, bookingID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM booking_payments WHERE booking_id = ?`, bookingID)
	return err
}

// GetByBookingID returns all payments for a specific booking.
func (r *BookingPaymentRepository) GetByBookingID(ctx context.Context, bookingID int) ([]models.BookingPayment, error) {
	query := `SELECT bp.id, bp.booking_id, bp.payment_type_id, bp.amount, pt.name
              FROM booking_payments bp
              LEFT JOIN payment_types pt ON bp.payment_type_id = pt.id
              WHERE bp.booking_id = ?`
	rows, err := r.db.QueryContext(ctx, query, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.BookingPayment
	for rows.Next() {
		var p models.BookingPayment
		if err := rows.Scan(&p.ID, &p.BookingID, &p.PaymentTypeID, &p.Amount, &p.PaymentType); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}
