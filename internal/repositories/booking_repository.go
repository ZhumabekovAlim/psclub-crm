package repositories

import (
	"context"
	"database/sql"
	"log"
	"psclub-crm/internal/models"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CreateWithItems(ctx context.Context, b *models.Booking) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin tx error: %v", err)
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `INSERT INTO bookings (client_id, table_id, user_id, start_time, end_time, note, discount, discount_reason, total_amount, bonus_used, payment_status, payment_type_id, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	res, err := tx.ExecContext(ctx, query, b.ClientID, b.TableID, b.UserID, b.StartTime, b.EndTime, b.Note, b.Discount, b.DiscountReason, b.TotalAmount, b.BonusUsed, b.PaymentStatus, b.PaymentTypeID)
	if err != nil {
		log.Printf("insert booking error: %v", err)
		return 0, err
	}
	bookingID, err := res.LastInsertId()
	if err != nil {
		log.Printf("last insert id error: %v", err)
		return 0, err
	}

	if len(b.Items) > 0 {
		itemQuery := `INSERT INTO booking_items (booking_id, item_id, quantity, price, discount) VALUES (?, ?, ?, ?, ?)`
		for _, item := range b.Items {
			_, err := tx.ExecContext(ctx, itemQuery, bookingID, item.ItemID, item.Quantity, item.Price, item.Discount)
			if err != nil {
				log.Printf("insert booking item error: %v", err)
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("commit booking error: %v", err)
		return 0, err
	}
	return int(bookingID), nil
}

func (r *BookingRepository) GetAll(ctx context.Context) ([]models.Booking, error) {
	query := `SELECT b.id, client_id, table_id, b.user_id, start_time, end_time, note, discount, discount_reason, total_amount, bonus_used, payment_status, payment_type_id, b.created_at, b.updated_at, c.name AS client_name, c.phone AS client_phone, payment_types.name AS payment_type, channels.name AS channel_name
				FROM bookings b
				JOIN clients c ON b.client_id = c.id
        		LEFT JOIN payment_types ON b.payment_type_id = payment_types.id
				LEFT JOIN channels ON c.channel_id = channels.id
				ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("get all bookings query error: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []models.Booking
	for rows.Next() {
		var b models.Booking
		err := rows.Scan(&b.ID, &b.ClientID, &b.TableID, &b.UserID, &b.StartTime, &b.EndTime, &b.Note, &b.Discount, &b.DiscountReason, &b.TotalAmount, &b.BonusUsed, &b.PaymentStatus, &b.PaymentTypeID, &b.CreatedAt, &b.UpdatedAt, &b.ClientName, &b.ClientPhone, &b.PaymentType, &b.ChannelName)
		if err != nil {
			log.Printf("scan booking error: %v", err)
			return nil, err
		}
		result = append(result, b)
	}
	return result, nil
}

// Получить все позиции по бронированию
func (r *BookingItemRepository) GetByBookingID(ctx context.Context, bookingID int) ([]models.BookingItem, error) {
	query := `SELECT bi.id, booking_id, item_id, bi.quantity, price, discount, pi.name FROM booking_items bi
                JOIN price_items pi ON bi.item_id = pi.id                                      
            	WHERE booking_id = ?`
	rows, err := r.db.QueryContext(ctx, query, bookingID)
	if err != nil {
		log.Printf("get booking items query error: %v", err)
		return nil, err
	}
	defer rows.Close()
	var items []models.BookingItem
	for rows.Next() {
		var it models.BookingItem
		err := rows.Scan(&it.ID, &it.BookingID, &it.ItemID, &it.Quantity, &it.Price, &it.Discount, &it.ItemName)
		if err != nil {
			log.Printf("scan booking item error: %v", err)
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	query := `SELECT bookings.id, client_id, table_id, user_id, start_time, end_time, note, discount, discount_reason, total_amount, bonus_used, payment_status, payment_type_id, bookings.created_at, bookings.updated_at, payment_types.name AS payment_type, channels.name AS channel_name
				FROM bookings 
        		LEFT JOIN payment_types ON bookings.payment_type_id = payment_types.id
				JOIN clients c ON bookings.client_id = c.id
				LEFT JOIN channels ON c.channel_id = channels.id
            	WHERE bookings.id = ?`
	var b models.Booking
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.ClientID, &b.TableID, &b.UserID, &b.StartTime, &b.EndTime, &b.Note, &b.Discount, &b.DiscountReason,
		&b.TotalAmount, &b.BonusUsed, &b.PaymentStatus, &b.PaymentTypeID, &b.CreatedAt, &b.UpdatedAt, &b.PaymentType, &b.ChannelName,
	)
	if err != nil {
		log.Printf("get booking by id error: %v", err)
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) Update(ctx context.Context, b *models.Booking) error {
	query := `UPDATE bookings SET client_id=?, table_id=?, user_id=?, start_time=?, end_time=?, note=?, discount=?, discount_reason=?, total_amount=?, bonus_used=?, payment_status=?, payment_type_id=?, updated_at=NOW() WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, b.ClientID, b.TableID, b.UserID, b.StartTime, b.EndTime, b.Note, b.Discount, b.DiscountReason, b.TotalAmount, b.BonusUsed, b.PaymentStatus, b.PaymentTypeID, b.ID)
	if err != nil {
		log.Printf("update booking error: %v", err)
	}
	return err
}

// UpdateWithItems updates booking data and replaces its items within a single transaction.
func (r *BookingRepository) UpdateWithItems(ctx context.Context, b *models.Booking) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("begin tx error: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `UPDATE bookings SET client_id=?, table_id=?, user_id=?, start_time=?, end_time=?, note=?, discount=?, discount_reason=?, total_amount=?, bonus_used=?, payment_status=?, payment_type_id=?, updated_at=NOW() WHERE id=?`
	_, err = tx.ExecContext(ctx, query, b.ClientID, b.TableID, b.UserID, b.StartTime, b.EndTime, b.Note, b.Discount, b.DiscountReason, b.TotalAmount, b.BonusUsed, b.PaymentStatus, b.PaymentTypeID, b.ID)
	if err != nil {
		log.Printf("update booking error: %v", err)
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM booking_items WHERE booking_id=?`, b.ID); err != nil {
		log.Printf("delete booking items error: %v", err)
		return err
	}

	if len(b.Items) > 0 {
		itemQuery := `INSERT INTO booking_items (booking_id, item_id, quantity, price, discount) VALUES (?, ?, ?, ?, ?)`
		for _, it := range b.Items {
			if _, err = tx.ExecContext(ctx, itemQuery, b.ID, it.ItemID, it.Quantity, it.Price, it.Discount); err != nil {
				log.Printf("insert booking item error: %v", err)
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("commit booking update error: %v", err)
	}
	return err
}

func (r *BookingRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM bookings WHERE id = ?`, id)
	if err != nil {
		log.Printf("delete booking error: %v", err)
	}
	return err
}
