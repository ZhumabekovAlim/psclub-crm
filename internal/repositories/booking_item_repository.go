package repositories

import (
	"database/sql"
)

type BookingItemRepository struct {
	db *sql.DB
}

func NewBookingItemRepository(db *sql.DB) *BookingItemRepository {
	return &BookingItemRepository{db: db}
}

//
//func (r *BookingItemRepository) GetByBookingID(ctx context.Context, bookingID int) ([]models.BookingItem, error) {
//	query := `SELECT id, booking_id, item_id, quantity, price, discount FROM booking_items WHERE booking_id = ?`
//	rows, err := r.db.QueryContext(ctx, query, bookingID)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	var items []models.BookingItem
//	for rows.Next() {
//		var it models.BookingItem
//		err := rows.Scan(&it.ID, &it.BookingID, &it.ItemID, &it.Quantity, &it.Price, &it.Discount)
//		if err != nil {
//			return nil, err
//		}
//		items = append(items, it)
//	}
//	return items, nil
//}
