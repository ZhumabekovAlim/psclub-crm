package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type PriceItemHistoryRepository struct {
	db *sql.DB
}

func NewPriceItemHistoryRepository(db *sql.DB) *PriceItemHistoryRepository {
	return &PriceItemHistoryRepository{db: db}
}

func (r *PriceItemHistoryRepository) Create(ctx context.Context, h *models.PriceItemHistory) (int, error) {
	query := `INSERT INTO price_item_history (price_item_id, operation, quantity, buy_price, total, user_id, created_at)
              VALUES (?, ?, ?, ?, ?, ?, NOW())`
	res, err := r.db.ExecContext(ctx, query, h.PriceItemID, h.Operation, h.Quantity, h.BuyPrice, h.Total, h.UserID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *PriceItemHistoryRepository) GetByItem(ctx context.Context, priceItemID int) ([]models.PriceItemHistory, error) {
	query := `SELECT id, price_item_id, operation, quantity, buy_price, total, user_id, created_at FROM price_item_history WHERE price_item_id = ? ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query, priceItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PriceItemHistory
	for rows.Next() {
		var h models.PriceItemHistory
		err := rows.Scan(&h.ID, &h.PriceItemID, &h.Operation, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (r *PriceItemHistoryRepository) GetAll(ctx context.Context) ([]models.PriceItemHistory, error) {
	query := `SELECT id, price_item_id, operation, quantity, buy_price, total, user_id, created_at FROM price_item_history ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PriceItemHistory
	for rows.Next() {
		var h models.PriceItemHistory
		err := rows.Scan(&h.ID, &h.PriceItemID, &h.Operation, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}
