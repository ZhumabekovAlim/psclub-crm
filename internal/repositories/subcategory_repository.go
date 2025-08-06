package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type SubcategoryRepository struct {
	db *sql.DB
}

func NewSubcategoryRepository(db *sql.DB) *SubcategoryRepository {
	return &SubcategoryRepository{db: db}
}

func (r *SubcategoryRepository) Create(ctx context.Context, s *models.Subcategory) (int, error) {
	query := `INSERT INTO subcategories (category_id, name, company_id, branch_id) VALUES (?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, s.CategoryID, s.Name, s.CompanyID, s.BranchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *SubcategoryRepository) GetAll(ctx context.Context, companyID, branchID int) ([]models.Subcategory, error) {
	query := `SELECT id, category_id, name, company_id, branch_id FROM subcategories WHERE company_id=? AND branch_id=? ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Subcategory
	for rows.Next() {
		var s models.Subcategory
		err := rows.Scan(&s.ID, &s.CategoryID, &s.Name, &s.CompanyID, &s.BranchID)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *SubcategoryRepository) GetByID(ctx context.Context, id, companyID, branchID int) (*models.Subcategory, error) {
	query := `SELECT id, category_id, name, company_id, branch_id FROM subcategories WHERE id=? AND company_id=? AND branch_id=?`
	var s models.Subcategory
	err := r.db.QueryRowContext(ctx, query, id, companyID, branchID).Scan(&s.ID, &s.CategoryID, &s.Name, &s.CompanyID, &s.BranchID)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubcategoryRepository) GetSubcategoriesByCategoryID(ctx context.Context, categoryID, companyID, branchID int) ([]models.Subcategory, error) {
	query := `SELECT id, category_id, name, company_id, branch_id FROM subcategories WHERE category_id = ? AND company_id=? AND branch_id=?`
	rows, err := r.db.QueryContext(ctx, query, categoryID, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Subcategory
	for rows.Next() {
		var s models.Subcategory
		if err := rows.Scan(&s.ID, &s.CategoryID, &s.Name, &s.CompanyID, &s.BranchID); err != nil {
			return nil, err
		}
		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *SubcategoryRepository) Update(ctx context.Context, s *models.Subcategory) error {
	query := `UPDATE subcategories SET category_id=?, name=? WHERE id=? AND company_id=? AND branch_id=?`
	_, err := r.db.ExecContext(ctx, query, s.CategoryID, s.Name, s.ID, s.CompanyID, s.BranchID)
	return err
}

func (r *SubcategoryRepository) Delete(ctx context.Context, id, companyID, branchID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM subcategories WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}
