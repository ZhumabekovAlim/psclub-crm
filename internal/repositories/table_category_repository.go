package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type TableCategoryRepository struct {
	db *sql.DB
}

func NewTableCategoryRepository(db *sql.DB) *TableCategoryRepository {
	return &TableCategoryRepository{db: db}
}

func (r *TableCategoryRepository) Create(ctx context.Context, c *models.TableCategory) (int, error) {
	query := `INSERT INTO table_categories (name, company_id, branch_id) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, c.Name, c.CompanyID, c.BranchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *TableCategoryRepository) GetAll(ctx context.Context, companyID, branchID int) ([]models.TableCategory, error) {
	query := `SELECT id, name, company_id, branch_id FROM table_categories WHERE company_id=? AND branch_id=? ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.TableCategory
	for rows.Next() {
		var c models.TableCategory
		err := rows.Scan(&c.ID, &c.Name, &c.CompanyID, &c.BranchID)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *TableCategoryRepository) GetByID(ctx context.Context, id, companyID, branchID int) (*models.TableCategory, error) {
	var category models.TableCategory
	err := r.db.QueryRowContext(ctx, "SELECT id, name, company_id, branch_id FROM table_categories WHERE id = ? AND company_id=? AND branch_id=?", id, companyID, branchID).
		Scan(&category.ID, &category.Name, &category.CompanyID, &category.BranchID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *TableCategoryRepository) Update(ctx context.Context, category *models.TableCategory) error {
	_, err := r.db.ExecContext(ctx, "UPDATE table_categories SET name = ? WHERE id = ? AND company_id=? AND branch_id=?", category.Name, category.ID, category.CompanyID, category.BranchID)
	return err
}

func (r *TableCategoryRepository) Delete(ctx context.Context, id, companyID, branchID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM table_categories WHERE id = ? AND company_id=? AND branch_id=?", id, companyID, branchID)
	return err
}
