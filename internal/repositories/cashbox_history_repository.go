package repositories

import (
	"context"
	"database/sql"
	"time"

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

func (r *CashboxHistoryRepository) GetByDate(ctx context.Context, date time.Time) ([]models.CashboxHistory, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, operation, amount, created_at FROM cashbox_history WHERE DATE(created_at)=? ORDER BY id`,
		date.Format("2006-01-02"))
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

func (r *CashboxHistoryRepository) GetByPeriod(ctx context.Context, start, end time.Time) ([]models.CashboxHistory, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, operation, amount, created_at FROM cashbox_history WHERE created_at >= ? AND created_at < ? ORDER BY id`,
		start, end)
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
