package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"psclub-crm/internal/models"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CreateWithItems(ctx context.Context, companyID, branchID int, b *models.Booking) (int, error) {
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

	query := `INSERT INTO bookings (company_id, branch_id, client_id, table_id, user_id, start_time, end_time, note, discount, discount_reason, total_amount, bonus_used, payment_status, payment_type_id, created_at, updated_at)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	var clientID interface{}
	if b.ClientID > 0 {
		clientID = b.ClientID
	} else {
		clientID = nil
	}

	var tableID interface{}
	if b.TableID > 0 {
		tableID = b.TableID
	} else {
		tableID = nil
	}

	res, err := tx.ExecContext(ctx, query, companyID, branchID, clientID, tableID, b.UserID, b.StartTime, b.EndTime, b.Note, b.Discount, b.DiscountReason, b.TotalAmount, b.BonusUsed, b.PaymentStatus, b.PaymentTypeID)
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
		itemQuery := `INSERT INTO booking_items (booking_id, company_id, branch_id, item_id, quantity, price, discount) VALUES (?, ?, ?, ?, ?, ?, ?)`
		for _, item := range b.Items {
			_, err := tx.ExecContext(ctx, itemQuery, bookingID, companyID, branchID, item.ItemID, item.Quantity, item.Price, item.Discount)
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

func (r *BookingRepository) GetAll(ctx context.Context, companyID, branchID int) ([]models.Booking, error) {
	query := `SELECT b.id, b.company_id, b.branch_id, b.client_id, table_id, b.user_id, start_time, end_time, note, discount, discount_reason, total_amount, bonus_used, payment_status, payment_type_id, b.created_at, b.updated_at,
                              IFNULL(c.name, ''), IFNULL(c.phone, ''), payment_types.name AS payment_type, IFNULL(channels.name, '') AS channel_name
                              FROM bookings b
                              LEFT JOIN clients c ON b.client_id = c.id
                              LEFT JOIN payment_types ON b.payment_type_id = payment_types.id
                              LEFT JOIN channels ON c.channel_id = channels.id
                              WHERE b.company_id = ? AND b.branch_id = ?
                              ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query, companyID, branchID)
	if err != nil {
		log.Printf("get all bookings query error: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []models.Booking
	for rows.Next() {
		var b models.Booking
		var clientID sql.NullInt64
		var tableID sql.NullInt64
		var channelName sql.NullString
		err := rows.Scan(&b.ID, &b.CompanyID, &b.BranchID, &clientID, &tableID, &b.UserID, &b.StartTime, &b.EndTime, &b.Note, &b.Discount, &b.DiscountReason, &b.TotalAmount, &b.BonusUsed, &b.PaymentStatus, &b.PaymentTypeID, &b.CreatedAt, &b.UpdatedAt, &b.ClientName, &b.ClientPhone, &b.PaymentType, &channelName)
		if err != nil {
			log.Printf("scan booking error: %v", err)
			return nil, err
		}
		if clientID.Valid {
			b.ClientID = int(clientID.Int64)
		}
		if tableID.Valid {
			b.TableID = int(tableID.Int64)
		}
		if channelName.Valid {
			b.ChannelName = channelName.String
		}
		result = append(result, b)
	}
	return result, nil
}

// GetByClientID returns all bookings for a specific client ordered by id DESC.
func (r *BookingRepository) GetByClientID(ctx context.Context, companyID, branchID, clientID int) ([]models.Booking, error) {
	query := `SELECT b.id, b.company_id, b.branch_id, b.client_id, table_id, b.user_id, start_time, end_time, note, discount, discount_reason, total_amount, bonus_used, payment_status, payment_type_id, b.created_at, b.updated_at,
                              IFNULL(c.name, ''), IFNULL(c.phone, ''), payment_types.name AS payment_type, IFNULL(channels.name, '') AS channel_name
                              FROM bookings b
                              LEFT JOIN clients c ON b.client_id = c.id
                              LEFT JOIN payment_types ON b.payment_type_id = payment_types.id
                              LEFT JOIN channels ON c.channel_id = channels.id
                              WHERE b.client_id = ? AND b.company_id = ? AND b.branch_id = ?
                              ORDER BY b.id DESC`
	rows, err := r.db.QueryContext(ctx, query, clientID, companyID, branchID)
	if err != nil {
		log.Printf("get bookings by client query error: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []models.Booking
	for rows.Next() {
		var b models.Booking
		var cID sql.NullInt64
		var tableID sql.NullInt64
		var channelName sql.NullString
		if err := rows.Scan(&b.ID, &b.CompanyID, &b.BranchID, &cID, &tableID, &b.UserID, &b.StartTime, &b.EndTime, &b.Note, &b.Discount, &b.DiscountReason, &b.TotalAmount, &b.BonusUsed, &b.PaymentStatus, &b.PaymentTypeID, &b.CreatedAt, &b.UpdatedAt, &b.ClientName, &b.ClientPhone, &b.PaymentType, &channelName); err != nil {
			log.Printf("scan booking by client error: %v", err)
			return nil, err
		}
		if cID.Valid {
			b.ClientID = int(cID.Int64)
		}
		if tableID.Valid {
			b.TableID = int(tableID.Int64)
		}
		if channelName.Valid {
			b.ChannelName = channelName.String
		}
		result = append(result, b)
	}
	return result, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, companyID, branchID, id int) (*models.Booking, error) {
	query := `SELECT bookings.id, bookings.company_id, bookings.branch_id, bookings.client_id, table_id, user_id, start_time, end_time, note, discount, discount_reason, total_amount, bonus_used, payment_status, payment_type_id, bookings.created_at, bookings.updated_at,
                              payment_types.name AS payment_type, IFNULL(channels.name, '') AS channel_name, IFNULL(c.name, ''), IFNULL(c.phone, '')
                              FROM bookings
                              LEFT JOIN payment_types ON bookings.payment_type_id = payment_types.id
                              LEFT JOIN clients c ON bookings.client_id = c.id
                              LEFT JOIN channels ON c.channel_id = channels.id
              WHERE bookings.id = ? AND bookings.company_id = ? AND bookings.branch_id = ?`
	var b models.Booking
	var clientID sql.NullInt64
	var tableID sql.NullInt64
	var channelName sql.NullString
	err := r.db.QueryRowContext(ctx, query, id, companyID, branchID).Scan(
		&b.ID, &b.CompanyID, &b.BranchID, &clientID, &tableID, &b.UserID, &b.StartTime, &b.EndTime, &b.Note, &b.Discount, &b.DiscountReason,
		&b.TotalAmount, &b.BonusUsed, &b.PaymentStatus, &b.PaymentTypeID, &b.CreatedAt, &b.UpdatedAt,
		&b.PaymentType, &channelName, &b.ClientName, &b.ClientPhone,
	)
	if err != nil {
		log.Printf("get booking by id error: %v", err)
		return nil, err
	}
	if clientID.Valid {
		b.ClientID = int(clientID.Int64)
	}
	if tableID.Valid {
		b.TableID = int(tableID.Int64)
	}
	if channelName.Valid {
		b.ChannelName = channelName.String
	}
	return &b, nil
}

func (r *BookingRepository) Update(ctx context.Context, companyID, branchID int, b *models.Booking) error {
	query := `UPDATE bookings SET company_id=?, branch_id=?, client_id=?, table_id=?, user_id=?, start_time=?, end_time=?, note=?, discount=?, discount_reason=?, total_amount=?, bonus_used=?, payment_status=?, payment_type_id=?, updated_at=NOW() WHERE id=? AND company_id=? AND branch_id=?`

	var clientID interface{}
	if b.ClientID > 0 {
		clientID = b.ClientID
	} else {
		clientID = nil
	}

	var tableID interface{}
	if b.TableID > 0 {
		tableID = b.TableID
	} else {
		tableID = nil
	}

	_, err := r.db.ExecContext(ctx, query, companyID, branchID, clientID, tableID, b.UserID, b.StartTime, b.EndTime, b.Note, b.Discount, b.DiscountReason, b.TotalAmount, b.BonusUsed, b.PaymentStatus, b.PaymentTypeID, b.ID, companyID, branchID)
	if err != nil {
		log.Printf("update booking error: %v", err)
	}
	return err
}

// UpdateWithItems updates booking data and replaces its items within a single transaction.
func (r *BookingRepository) UpdateWithItems(ctx context.Context, companyID, branchID int, b *models.Booking) (err error) {
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

	query := `UPDATE bookings SET company_id=?, branch_id=?, client_id=?, table_id=?, user_id=?, start_time=?, end_time=?, note=?, discount=?, discount_reason=?, total_amount=?, bonus_used=?, payment_status=?, payment_type_id=?, updated_at=NOW() WHERE id=? AND company_id=? AND branch_id=?`

	var clientID interface{}
	if b.ClientID > 0 {
		clientID = b.ClientID
	} else {
		clientID = nil
	}

	var tableID interface{}
	if b.TableID > 0 {
		tableID = b.TableID
	} else {
		tableID = nil
	}

	_, err = tx.ExecContext(ctx, query, companyID, branchID, clientID, tableID, b.UserID, b.StartTime, b.EndTime, b.Note, b.Discount, b.DiscountReason, b.TotalAmount, b.BonusUsed, b.PaymentStatus, b.PaymentTypeID, b.ID, companyID, branchID)
	if err != nil {
		log.Printf("update booking error: %v", err)
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM booking_items WHERE booking_id=? AND company_id=? AND branch_id=?`, b.ID, companyID, branchID); err != nil {
		log.Printf("delete booking items error: %v", err)
		return err
	}

	if len(b.Items) > 0 {
		itemQuery := `INSERT INTO booking_items (booking_id, company_id, branch_id, item_id, quantity, price, discount) VALUES (?, ?, ?, ?, ?, ?, ?)`
		for _, it := range b.Items {
			if _, err = tx.ExecContext(ctx, itemQuery, b.ID, companyID, branchID, it.ItemID, it.Quantity, it.Price, it.Discount); err != nil {
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

func (r *BookingRepository) Delete(ctx context.Context, companyID, branchID, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM bookings WHERE id = ? AND company_id = ? AND branch_id = ?`, id, companyID, branchID)
	if err != nil {
		log.Printf("delete booking error: %v", err)
	}
	return err
}

func (r *BookingRepository) IsTableAvailable(ctx context.Context, companyID, branchID, tableID, excludeID int, start, end time.Time) (bool, error) {
	query := `SELECT COUNT(1) FROM bookings WHERE company_id=? AND branch_id=? AND table_id=? AND ? < end_time AND ? > start_time`
	args := []interface{}{companyID, branchID, tableID, start, end}
	if excludeID > 0 {
		query += " AND id <> ?"
		args = append(args, excludeID)
	}
	var cnt int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&cnt)
	if err != nil {
		log.Printf("check table availability error: %v", err)
		return false, err
	}
	return cnt == 0, nil
}
