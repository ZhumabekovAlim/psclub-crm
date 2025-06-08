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

func (r *CategoryRepository) Create(ctx context.Context, c *models.Category) (int, error) {
	query := `INSERT INTO categories (name) VALUES (?)`
	res, err := r.db.ExecContext(ctx, query, c.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	query := `SELECT id, name FROM categories ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name)
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

func (r *CategoryRepository) Update(ctx context.Context, c *models.Category) error {
	query := `UPDATE categories SET name=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, c.Name, c.ID)
	return err
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id=?`, id)
	return err
}
