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
	query := `INSERT INTO subcategories (category_id, name) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, s.CategoryID, s.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *SubcategoryRepository) GetAll(ctx context.Context) ([]models.Subcategory, error) {
	query := `SELECT id, category_id, name FROM subcategories ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Subcategory
	for rows.Next() {
		var s models.Subcategory
		err := rows.Scan(&s.ID, &s.CategoryID, &s.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *SubcategoryRepository) GetByID(ctx context.Context, id int) (*models.Subcategory, error) {
	query := `SELECT id, category_id, name FROM subcategories WHERE id=?`
	var s models.Subcategory
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.CategoryID, &s.Name)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubcategoryRepository) GetSubcategoryByCategoryID(ctx context.Context, id int) (*models.Subcategory, error) {
	query := `SELECT id, category_id, name FROM subcategories WHERE subcategories.category_id=?`
	var s models.Subcategory
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.CategoryID, &s.Name)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubcategoryRepository) Update(ctx context.Context, s *models.Subcategory) error {
	query := `UPDATE subcategories SET category_id=?, name=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, s.CategoryID, s.Name, s.ID)
	return err
}

func (r *SubcategoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM subcategories WHERE id=?`, id)
	return err
}
