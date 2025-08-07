package repositories

import (
	"context"
	"database/sql"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
)

type CashboxRepository struct {
	db *sql.DB
}

func NewCashboxRepository(db *sql.DB) *CashboxRepository {
	return &CashboxRepository{db: db}
}

func (r *CashboxRepository) Get(ctx context.Context) (*models.Cashbox, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT id, amount FROM cashbox WHERE company_id=? AND branch_id=? LIMIT 1`
	var c models.Cashbox
	err := r.db.QueryRowContext(ctx, query, companyID, branchID).Scan(&c.ID, &c.Amount)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CashboxRepository) Update(ctx context.Context, c *models.Cashbox) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `UPDATE cashbox SET amount=? WHERE id=? AND company_id=? AND branch_id=?`
	_, err := r.db.ExecContext(ctx, query, c.Amount, c.ID, companyID, branchID)
	return err
}
