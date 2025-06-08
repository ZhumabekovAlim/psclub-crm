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
	query := `INSERT INTO table_categories (name) VALUES (?)`
	res, err := r.db.ExecContext(ctx, query, c.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *TableCategoryRepository) GetAll(ctx context.Context) ([]models.TableCategory, error) {
	query := `SELECT id, name FROM table_categories ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.TableCategory
	for rows.Next() {
		var c models.TableCategory
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *TableCategoryRepository) GetByID(id int) (*models.TableCategory, error) {
	var category models.TableCategory
	err := r.db.QueryRow("SELECT id, name FROM table_categories WHERE id = ?", id).
		Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *TableCategoryRepository) Update(id int, category *models.TableCategory) error {
	_, err := r.db.Exec("UPDATE table_categories SET name = ? WHERE id = ?", category.Name, id)
	return err
}

func (r *TableCategoryRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM table_categories WHERE id = ?", id)
	return err
}
