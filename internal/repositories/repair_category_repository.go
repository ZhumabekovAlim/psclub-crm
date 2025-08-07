package repositories

import (
	"context"
	"database/sql"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
)

type RepairCategoryRepository struct {
	db *sql.DB
}

func NewRepairCategoryRepository(db *sql.DB) *RepairCategoryRepository {
	return &RepairCategoryRepository{db: db}
}

func (r *RepairCategoryRepository) Create(ctx context.Context, c *models.RepairCategory) (int, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	res, err := r.db.ExecContext(ctx, `INSERT INTO repair_categories (name, company_id, branch_id) VALUES (?, ?, ?)`, c.Name, companyID, branchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *RepairCategoryRepository) GetAll(ctx context.Context) ([]models.RepairCategory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, company_id, branch_id FROM repair_categories WHERE company_id=? AND branch_id=? ORDER BY id`, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.RepairCategory
	for rows.Next() {
		var c models.RepairCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.CompanyID, &c.BranchID); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *RepairCategoryRepository) Update(ctx context.Context, c *models.RepairCategory) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `UPDATE repair_categories SET name=? WHERE id=? AND company_id=? AND branch_id=?`, c.Name, c.ID, companyID, branchID)
	return err
}

func (r *RepairCategoryRepository) Delete(ctx context.Context, id int) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `DELETE FROM repair_categories WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}

func (r *RepairCategoryRepository) GetByName(ctx context.Context, name string) (*models.RepairCategory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	var c models.RepairCategory
	err := r.db.QueryRowContext(ctx, `SELECT id, name, company_id, branch_id FROM repair_categories WHERE name=? AND company_id=? AND branch_id=?`, name, companyID, branchID).Scan(&c.ID, &c.Name, &c.CompanyID, &c.BranchID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
