package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type RepairRepository struct {
	db *sql.DB
}

func NewRepairRepository(db *sql.DB) *RepairRepository {
	return &RepairRepository{db: db}
}

func (r *RepairRepository) Create(ctx context.Context, rep *models.Repair) (int, error) {
	query := `INSERT INTO repairs (date, vin, description, price, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query, rep.Date, rep.VIN, rep.Description, rep.Price)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *RepairRepository) GetAll(ctx context.Context) ([]models.Repair, error) {
	query := `SELECT id, date, vin, description, price, created_at, updated_at FROM repairs ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Repair
	for rows.Next() {
		var rep models.Repair
		err := rows.Scan(&rep.ID, &rep.Date, &rep.VIN, &rep.Description, &rep.Price, &rep.CreatedAt, &rep.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, rep)
	}
	return result, nil
}

func (r *RepairRepository) GetByID(ctx context.Context, id int) (*models.Repair, error) {
	query := `SELECT id, date, vin, description, price, created_at, updated_at FROM repairs WHERE id = ?`
	var rep models.Repair
	err := r.db.QueryRowContext(ctx, query, id).Scan(&rep.ID, &rep.Date, &rep.VIN, &rep.Description, &rep.Price, &rep.CreatedAt, &rep.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rep, nil
}

func (r *RepairRepository) Update(ctx context.Context, rep *models.Repair) error {
	query := `UPDATE repairs SET date = ?, vin = ?, description = ?, price = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, rep.Date, rep.VIN, rep.Description, rep.Price, rep.ID)
	return err
}

func (r *RepairRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM repairs WHERE id = ?`, id)
	return err
}
