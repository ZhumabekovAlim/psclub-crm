package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type CashboxHistoryRepository struct {
	db *sql.DB
}

func NewCashboxHistoryRepository(db *sql.DB) *CashboxHistoryRepository {
	return &CashboxHistoryRepository{db: db}
}

func (r *CashboxHistoryRepository) Create(ctx context.Context, h *models.CashboxHistory) (int, error) {
	query := `INSERT INTO cashbox_history (operation, amount, created_at) VALUES (?, ?, NOW())`
	res, err := r.db.ExecContext(ctx, query, h.Operation, h.Amount)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *CashboxHistoryRepository) GetAll(ctx context.Context) ([]models.CashboxHistory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, operation, amount, created_at FROM cashbox_history ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.CashboxHistory
	for rows.Next() {
		var h models.CashboxHistory
		if err := rows.Scan(&h.ID, &h.Operation, &h.Amount, &h.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, nil
}
