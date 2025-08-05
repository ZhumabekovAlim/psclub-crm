package repositories

import (
	"context"
	"database/sql"

	"psclub-crm/internal/models"
)

type BookingItemRepository struct {
	db *sql.DB
}

func NewBookingItemRepository(db *sql.DB) *BookingItemRepository {
	return &BookingItemRepository{db: db}
}

// GetByBookingID returns booking items for the specified booking filtered by company and branch.
func (r *BookingItemRepository) GetByBookingID(ctx context.Context, companyID, branchID, bookingID int) ([]models.BookingItem, error) {
	query := `SELECT bi.id, bi.booking_id, bi.company_id, bi.branch_id, bi.item_id, bi.quantity, bi.price, bi.discount, pi.name
               FROM booking_items bi
               JOIN price_items pi ON bi.item_id = pi.id
               WHERE bi.booking_id = ? AND bi.company_id = ? AND bi.branch_id = ?`
	rows, err := r.db.QueryContext(ctx, query, bookingID, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.BookingItem
	for rows.Next() {
		var it models.BookingItem
		if err := rows.Scan(&it.ID, &it.BookingID, &it.CompanyID, &it.BranchID, &it.ItemID, &it.Quantity, &it.Price, &it.Discount, &it.ItemName); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}
