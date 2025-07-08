package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
	"time"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(ctx context.Context, e *models.Expense) (int, error) {
	query := `INSERT INTO expenses (date, title, category_id, repair_category_id, total, description, paid, created_at)
                VALUES (?, ?, NULLIF(?,0), NULLIF(?,0), ?, ?, ?, NOW())`
	res, err := r.db.ExecContext(ctx, query, e.Date, e.Title, e.CategoryID, e.RepairCategoryID, e.Total, e.Description, e.Paid)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ExpenseRepository) GetAll(ctx context.Context, from, to time.Time) ([]models.Expense, error) {
	query := `SELECT e.id, e.date, e.title, IFNULL(e.category_id, 0), IFNULL(ec.name, ''), IFNULL(e.repair_category_id, 0), IFNULL(rc.name,''), e.total, e.description, e.paid, e.created_at
                FROM expenses e
                LEFT JOIN expense_categories ec ON e.category_id = ec.id
                LEFT JOIN repair_categories rc ON e.repair_category_id = rc.id
                WHERE DATE(e.date) BETWEEN ? AND ?
                ORDER BY e.id DESC`
	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Date, &e.Title, &e.CategoryID, &e.Category, &e.RepairCategoryID, &e.RepairCategory, &e.Total, &e.Description, &e.Paid, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *ExpenseRepository) GetByID(ctx context.Context, id int) (*models.Expense, error) {
	query := `SELECT e.id, e.date, e.title, IFNULL(e.category_id, 0), IFNULL(ec.name, ''), IFNULL(e.repair_category_id, 0), IFNULL(rc.name,''), e.total, e.description, e.paid, e.created_at
                FROM expenses e
                LEFT JOIN expense_categories ec ON e.category_id = ec.id
                LEFT JOIN repair_categories rc ON e.repair_category_id = rc.id
                WHERE e.id = ?`
	var e models.Expense
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Date, &e.Title, &e.CategoryID, &e.Category, &e.RepairCategoryID, &e.RepairCategory, &e.Total, &e.Description, &e.Paid, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *ExpenseRepository) Update(ctx context.Context, e *models.Expense) error {
	query := `UPDATE expenses SET date=?, title=?, category_id=NULLIF(?,0), repair_category_id=NULLIF(?,0), total=?, description=?, paid=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, e.Date, e.Title, e.CategoryID, e.RepairCategoryID, e.Total, e.Description, e.Paid, e.ID)
	return err
}

func (r *ExpenseRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE id=?`, id)
	return err
}

func (r *ExpenseRepository) DeleteByDetails(ctx context.Context, title, description string, total float64, repairCategoryID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE title=? AND description=? AND total=? AND IFNULL(repair_category_id,0)=?`, title, description, total, repairCategoryID)
	return err
}
