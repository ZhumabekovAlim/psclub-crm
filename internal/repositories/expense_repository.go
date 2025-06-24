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
	query := `INSERT INTO expenses (date, title, category_id, total, description, paid, created_at)
                VALUES (?, ?, NULLIF(?,0), ?, ?, ?, NOW())`
	res, err := r.db.ExecContext(ctx, query, e.Date, e.Title, e.CategoryID, e.Total, e.Description, e.Paid)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ExpenseRepository) GetAll(ctx context.Context) ([]models.Expense, error) {
	query := `SELECT e.id, e.date, e.title, e.category_id, IFNULL(ec.name, ''), e.total, e.description, e.paid, e.created_at
                FROM expenses e
                LEFT JOIN expense_categories ec ON e.category_id = ec.id
                ORDER BY e.id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Date, &e.Title, &e.CategoryID, &e.Category, &e.Total, &e.Description, &e.Paid, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *ExpenseRepository) GetByID(ctx context.Context, id int) (*models.Expense, error) {
	query := `SELECT e.id, e.date, e.title, e.category_id, IFNULL(ec.name, ''), e.total, e.description, e.paid, e.created_at
                FROM expenses e
                LEFT JOIN expense_categories ec ON e.category_id = ec.id
                WHERE e.id = ?`
	var e models.Expense
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Date, &e.Title, &e.CategoryID, &e.Category, &e.Total, &e.Description, &e.Paid, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *ExpenseRepository) Update(ctx context.Context, e *models.Expense) error {
	query := `UPDATE expenses SET date=?, title=?, category_id=NULLIF(?,0), total=?, description=?, paid=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, e.Date, e.Title, e.CategoryID, e.Total, e.Description, e.Paid, e.ID)
	return err
}

func (r *ExpenseRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE id=?`, id)
	return err
}

func (r *ExpenseRepository) DeleteByDetails(ctx context.Context, title, description string, total float64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE title=? AND description=? AND total=?`, title, description, total)
	return err
}
