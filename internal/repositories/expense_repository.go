package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(ctx context.Context, e *models.Expense) (int, error) {
	query := `INSERT INTO expenses (date, title, category, total, description, created_at) VALUES (?, ?, ?, ?, ?, NOW())`
	res, err := r.db.ExecContext(ctx, query, e.Date, e.Title, e.Category, e.Total, e.Description)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ExpenseRepository) GetAll(ctx context.Context) ([]models.Expense, error) {
	query := `SELECT id, date, title, category, total, description, created_at FROM expenses ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Date, &e.Title, &e.Category, &e.Total, &e.Description, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *ExpenseRepository) GetByID(ctx context.Context, id int) (*models.Expense, error) {
	query := `SELECT id, date, title, category, total, description, created_at FROM expenses WHERE id = ?`
	var e models.Expense
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Date, &e.Title, &e.Category, &e.Total, &e.Description, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *ExpenseRepository) Update(ctx context.Context, e *models.Expense) error {
	query := `UPDATE expenses SET date=?, title=?, category=?, total=?, description=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, e.Date, e.Title, e.Category, e.Total, e.Description, e.ID)
	return err
}

func (r *ExpenseRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE id=?`, id)
	return err
}
