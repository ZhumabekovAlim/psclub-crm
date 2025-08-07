package repositories

import (
	"context"
	"database/sql"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
)

type ExpenseCategoryRepository struct {
	db *sql.DB
}

func NewExpenseCategoryRepository(db *sql.DB) *ExpenseCategoryRepository {
	return &ExpenseCategoryRepository{db: db}
}

func (r *ExpenseCategoryRepository) Create(ctx context.Context, c *models.ExpenseCategory) (int, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	res, err := r.db.ExecContext(ctx, `INSERT INTO expense_categories (name, company_id, branch_id) VALUES (?, ?, ?)`, c.Name, companyID, branchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ExpenseCategoryRepository) GetAll(ctx context.Context) ([]models.ExpenseCategory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, company_id, branch_id FROM expense_categories WHERE company_id=? AND branch_id=? ORDER BY id`, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ExpenseCategory
	for rows.Next() {
		var c models.ExpenseCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.CompanyID, &c.BranchID); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *ExpenseCategoryRepository) Update(ctx context.Context, c *models.ExpenseCategory) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `UPDATE expense_categories SET name=? WHERE id=? AND company_id=? AND branch_id=?`, c.Name, c.ID, companyID, branchID)
	return err
}

func (r *ExpenseCategoryRepository) Delete(ctx context.Context, id int) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `DELETE FROM expense_categories WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}

func (r *ExpenseCategoryRepository) GetByName(ctx context.Context, name string) (*models.ExpenseCategory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT id, name, company_id, branch_id FROM expense_categories WHERE name = ? AND company_id=? AND branch_id=?`
	var c models.ExpenseCategory
	err := r.db.QueryRowContext(ctx, query, name, companyID, branchID).Scan(&c.ID, &c.Name, &c.CompanyID, &c.BranchID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
