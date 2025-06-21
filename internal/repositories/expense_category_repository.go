package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type ExpenseCategoryRepository struct {
	db *sql.DB
}

func NewExpenseCategoryRepository(db *sql.DB) *ExpenseCategoryRepository {
	return &ExpenseCategoryRepository{db: db}
}

func (r *ExpenseCategoryRepository) Create(ctx context.Context, c *models.ExpenseCategory) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO expense_categories (name) VALUES (?)`, c.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ExpenseCategoryRepository) GetAll(ctx context.Context) ([]models.ExpenseCategory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM expense_categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ExpenseCategory
	for rows.Next() {
		var c models.ExpenseCategory
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *ExpenseCategoryRepository) Update(ctx context.Context, c *models.ExpenseCategory) error {
	_, err := r.db.ExecContext(ctx, `UPDATE expense_categories SET name=? WHERE id=?`, c.Name, c.ID)
	return err
}

func (r *ExpenseCategoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expense_categories WHERE id=?`, id)
	return err
}

func (r *ExpenseCategoryRepository) GetByName(ctx context.Context, name string) (*models.ExpenseCategory, error) {
	query := `SELECT id, name FROM expense_categories WHERE name = ?`
	var c models.ExpenseCategory
	err := r.db.QueryRowContext(ctx, query, name).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
