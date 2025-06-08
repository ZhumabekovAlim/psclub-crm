package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type CashboxRepository struct {
	db *sql.DB
}

func NewCashboxRepository(db *sql.DB) *CashboxRepository {
	return &CashboxRepository{db: db}
}

func (r *CashboxRepository) Get(ctx context.Context) (*models.Cashbox, error) {
	query := `SELECT id, amount FROM cashbox LIMIT 1`
	var c models.Cashbox
	err := r.db.QueryRowContext(ctx, query).Scan(&c.ID, &c.Amount)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CashboxRepository) Update(ctx context.Context, c *models.Cashbox) error {
	query := `UPDATE cashbox SET amount=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, c.Amount, c.ID)
	return err
}
