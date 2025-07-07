package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type RepairCategoryRepository struct {
	db *sql.DB
}

func NewRepairCategoryRepository(db *sql.DB) *RepairCategoryRepository {
	return &RepairCategoryRepository{db: db}
}

func (r *RepairCategoryRepository) Create(ctx context.Context, c *models.RepairCategory) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO repair_categories (name) VALUES (?)`, c.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *RepairCategoryRepository) GetAll(ctx context.Context) ([]models.RepairCategory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM repair_categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.RepairCategory
	for rows.Next() {
		var c models.RepairCategory
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *RepairCategoryRepository) Update(ctx context.Context, c *models.RepairCategory) error {
	_, err := r.db.ExecContext(ctx, `UPDATE repair_categories SET name=? WHERE id=?`, c.Name, c.ID)
	return err
}

func (r *RepairCategoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM repair_categories WHERE id=?`, id)
	return err
}

func (r *RepairCategoryRepository) GetByName(ctx context.Context, name string) (*models.RepairCategory, error) {
	var c models.RepairCategory
	err := r.db.QueryRowContext(ctx, `SELECT id, name FROM repair_categories WHERE name=?`, name).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
