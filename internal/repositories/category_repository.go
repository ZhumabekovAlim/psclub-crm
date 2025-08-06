package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetByName(ctx context.Context, name string, companyID, branchID int) (*models.Category, error) {
	query := `SELECT id, name, company_id, branch_id FROM categories WHERE name = ? AND company_id = ? AND branch_id = ?`
	var c models.Category
	err := r.db.QueryRowContext(ctx, query, name, companyID, branchID).Scan(&c.ID, &c.Name, &c.CompanyID, &c.BranchID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) (int, error) {
	query := `INSERT INTO categories (name, company_id, branch_id) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, c.Name, c.CompanyID, c.BranchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *CategoryRepository) GetAll(ctx context.Context, companyID, branchID int) ([]models.Category, error) {
	query := `SELECT id, name, company_id, branch_id FROM categories WHERE company_id=? AND branch_id=? ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.CompanyID, &c.BranchID)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	query := `SELECT id, name FROM categories WHERE id=?`
	var c models.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) GetByIDTenant(ctx context.Context, id, companyID, branchID int) (*models.Category, error) {
	query := `SELECT id, name, company_id, branch_id FROM categories WHERE id=? AND company_id=? AND branch_id=?`
	var c models.Category
	err := r.db.QueryRowContext(ctx, query, id, companyID, branchID).Scan(&c.ID, &c.Name, &c.CompanyID, &c.BranchID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *models.Category) error {
	query := `UPDATE categories SET name=? WHERE id=? AND company_id=? AND branch_id=?`
	_, err := r.db.ExecContext(ctx, query, c.Name, c.ID, c.CompanyID, c.BranchID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id, companyID, branchID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}
