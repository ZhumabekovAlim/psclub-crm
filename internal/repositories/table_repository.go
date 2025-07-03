package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) Create(ctx context.Context, t *models.Table) (int, error) {
	// Determine the next table number automatically
	err := r.db.QueryRowContext(ctx, `SELECT IFNULL(MAX(number), 0) + 1 FROM tables`).Scan(&t.Number)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO tables (category_id, name, number) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, t.CategoryID, t.Name, t.Number)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *TableRepository) GetAll(ctx context.Context) ([]models.Table, error) {
	query := `SELECT id, category_id, name, number FROM tables ORDER BY number`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Table
	for rows.Next() {
		var t models.Table
		err := rows.Scan(&t.ID, &t.CategoryID, &t.Name, &t.Number)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *TableRepository) GetByID(id int) (*models.Table, error) {
	var table models.Table
	err := r.db.QueryRow("SELECT id, name, category_id, number FROM tables WHERE id = ?", id).
		Scan(&table.ID, &table.Name, &table.CategoryID, &table.Number)
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *TableRepository) Update(id int, table *models.Table) error {
	_, err := r.db.Exec("UPDATE tables SET name = ?, category_id = ?, number = ? WHERE id = ?", table.Name, table.CategoryID, table.Number, id)
	return err
}

func (r *TableRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM tables WHERE id = ?", id)
	return err
}

func (r *TableRepository) Reorder(ctx context.Context, id int, newNumber int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var current int
	if err := tx.QueryRowContext(ctx, `SELECT number FROM tables WHERE id = ?`, id).Scan(&current); err != nil {
		tx.Rollback()
		return err
	}
	if newNumber == current {
		return tx.Commit()
	}
	if newNumber < current {
		if _, err := tx.ExecContext(ctx, `UPDATE tables SET number = number + 1 WHERE number >= ? AND number < ?`, newNumber, current); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if _, err := tx.ExecContext(ctx, `UPDATE tables SET number = number - 1 WHERE number > ? AND number <= ?`, current, newNumber); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE tables SET number = ? WHERE id = ?`, newNumber, id); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
