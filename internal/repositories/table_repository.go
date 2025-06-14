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
	query := `INSERT INTO tables (category_id, name) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, t.CategoryID, t.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *TableRepository) GetAll(ctx context.Context) ([]models.Table, error) {
	query := `SELECT id, category_id, name FROM tables ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Table
	for rows.Next() {
		var t models.Table
		err := rows.Scan(&t.ID, &t.CategoryID, &t.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *TableRepository) GetByID(id int) (*models.Table, error) {
	var table models.Table
	err := r.db.QueryRow("SELECT id, name, category_id FROM tables WHERE id = ?", id).
		Scan(&table.ID, &table.Name, &table.CategoryID)
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *TableRepository) Update(id int, table *models.Table) error {
	_, err := r.db.Exec("UPDATE tables SET name = ?, category_id = ? WHERE id = ?", table.Name, table.CategoryID, id)
	return err
}

func (r *TableRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM tables WHERE id = ?", id)
	return err
}
