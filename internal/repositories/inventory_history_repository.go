package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type InventoryHistoryRepository struct {
	db *sql.DB
}

func NewInventoryHistoryRepository(db *sql.DB) *InventoryHistoryRepository {
	return &InventoryHistoryRepository{db: db}
}

func (r *InventoryHistoryRepository) Create(ctx context.Context, h *models.InventoryHistory) (int, error) {
	query := `INSERT INTO inventory_history (price_item_id, expected, actual, difference, created_at) VALUES (?, ?, ?, ?, NOW())`
	res, err := r.db.ExecContext(ctx, query, h.PriceItemID, h.Expected, h.Actual, h.Difference)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *InventoryHistoryRepository) GetAll(ctx context.Context, companyID, branchID int) ([]models.InventoryHistory, error) {
	rows, err := r.db.QueryContext(ctx, `
			SELECT ih.id, ih.price_item_id, pi.name, ih.expected, ih.actual, ih.difference, ih.created_at
			FROM inventory_history ih
			LEFT JOIN price_items pi ON ih.price_item_id = pi.id
			WHERE ih.company_id = ? AND ih.branch_id = ?
			ORDER BY ih.id DESC`, companyID, branchID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.InventoryHistory
	for rows.Next() {
		var h models.InventoryHistory
		if err := rows.Scan(&h.ID, &h.PriceItemID, &h.Name, &h.Expected, &h.Actual, &h.Difference, &h.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, nil
}
