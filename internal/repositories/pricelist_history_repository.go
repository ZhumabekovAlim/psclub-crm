package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type PricelistHistoryRepository struct {
	db *sql.DB
}

func NewPricelistHistoryRepository(db *sql.DB) *PricelistHistoryRepository {
	return &PricelistHistoryRepository{db: db}
}

func (r *PricelistHistoryRepository) Create(ctx context.Context, h *models.PricelistHistory) (int, error) {
	query := `INSERT INTO pricelist_history (price_item_id, quantity, buy_price, total, user_id, created_at) VALUES (?, ?, ?, ?, ?, NOW())`
	res, err := r.db.ExecContext(ctx, query, h.PriceItemID, h.Quantity, h.BuyPrice, h.Total, h.UserID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *PricelistHistoryRepository) GetByItem(ctx context.Context, priceItemID int) ([]models.PricelistHistory, error) {
	query := `SELECT id, price_item_id, quantity, buy_price, total, user_id, created_at FROM pricelist_history WHERE price_item_id = ? ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query, priceItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PricelistHistory
	for rows.Next() {
		var h models.PricelistHistory
		if err := rows.Scan(&h.ID, &h.PriceItemID, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (r *PricelistHistoryRepository) GetAll(ctx context.Context) ([]models.PricelistHistory, error) {
	query := `SELECT id, price_item_id, quantity, buy_price, total, user_id, created_at FROM pricelist_history ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PricelistHistory
	for rows.Next() {
		var h models.PricelistHistory
		if err := rows.Scan(&h.ID, &h.PriceItemID, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}
func (r *PricelistHistoryRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.PricelistHistory, error) {
	query := `SELECT ph.id, ph.price_item_id, ph.quantity, ph.buy_price, ph.total, ph.user_id, ph.created_at, u.name AS user_name
                FROM pricelist_history ph
                JOIN price_items pi ON ph.price_item_id = pi.id
                JOIN users u ON ph.user_id = u.id
                WHERE pi.category_id = ? ORDER BY ph.id DESC`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PricelistHistory
	for rows.Next() {
		var h models.PricelistHistory
		if err := rows.Scan(&h.ID, &h.PriceItemID, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CreatedAt, &h.UserName); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}
