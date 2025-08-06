package repositories

import (
	"context"
	"database/sql"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
)

// PaymentTypeRepository handles CRUD operations for payment types.
type PaymentTypeRepository struct {
	db *sql.DB
}

func NewPaymentTypeRepository(db *sql.DB) *PaymentTypeRepository {
	return &PaymentTypeRepository{db: db}
}

func (r *PaymentTypeRepository) GetByID(ctx context.Context, id int) (*models.PaymentType, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	var pt models.PaymentType
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, hold_percent, company_id, branch_id FROM payment_types WHERE id=? AND company_id=? AND branch_id=?`,
		id, companyID, branchID,
	).Scan(&pt.ID, &pt.Name, &pt.HoldPercent, &pt.CompanyID, &pt.BranchID)
	if err != nil {
		return nil, err
	}
	return &pt, nil
}

func (r *PaymentTypeRepository) GetAll(ctx context.Context) ([]models.PaymentType, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, hold_percent, company_id, branch_id FROM payment_types WHERE company_id=? AND branch_id=? ORDER BY id`, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.PaymentType
	for rows.Next() {
		var pt models.PaymentType
		if err := rows.Scan(&pt.ID, &pt.Name, &pt.HoldPercent, &pt.CompanyID, &pt.BranchID); err != nil {
			return nil, err
		}
		result = append(result, pt)
	}
	return result, nil
}

func (r *PaymentTypeRepository) Create(ctx context.Context, pt *models.PaymentType) (int, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	res, err := r.db.ExecContext(ctx, `INSERT INTO payment_types (name, hold_percent, company_id, branch_id) VALUES (?, ?, ?, ?)`, pt.Name, pt.HoldPercent, companyID, branchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *PaymentTypeRepository) Update(ctx context.Context, pt *models.PaymentType) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx,
		`UPDATE payment_types SET name=?, hold_percent=? WHERE id=? AND company_id=? AND branch_id=?`,
		pt.Name, pt.HoldPercent, pt.ID, companyID, branchID,
	)
	return err
}

func (r *PaymentTypeRepository) Delete(ctx context.Context, id int) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `DELETE FROM payment_types WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}
