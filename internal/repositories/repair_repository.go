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
	query := `INSERT INTO repairs (date, vin, description, price, category_id, created_at, updated_at)
                VALUES (?, ?, ?, ?, NULLIF(?,0), NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query, rep.Date, rep.VIN, rep.Description, rep.Price, rep.CategoryID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *RepairRepository) GetAll(ctx context.Context) ([]models.Repair, error) {
	query := `SELECT r.id, r.date, r.vin, r.description, r.price, r.category_id, IFNULL(rc.name, ''), r.created_at, r.updated_at
                FROM repairs r
                LEFT JOIN repair_categories rc ON r.category_id = rc.id
                ORDER BY r.id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Repair
	for rows.Next() {
		var rep models.Repair
		err := rows.Scan(&rep.ID, &rep.Date, &rep.VIN, &rep.Description, &rep.Price, &rep.CategoryID, &rep.Category, &rep.CreatedAt, &rep.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, rep)
	}
	return result, nil
}

func (r *RepairRepository) GetByID(ctx context.Context, id int) (*models.Repair, error) {
	query := `SELECT r.id, r.date, r.vin, r.description, r.price, r.category_id, IFNULL(rc.name,''), r.created_at, r.updated_at
                FROM repairs r
                LEFT JOIN repair_categories rc ON r.category_id = rc.id
                WHERE r.id = ?`
	var rep models.Repair
	err := r.db.QueryRowContext(ctx, query, id).Scan(&rep.ID, &rep.Date, &rep.VIN, &rep.Description, &rep.Price, &rep.CategoryID, &rep.Category, &rep.CreatedAt, &rep.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &rep, nil
}

func (r *RepairRepository) Update(ctx context.Context, rep *models.Repair) error {
	query := `UPDATE repairs SET date = ?, vin = ?, description = ?, price = ?, category_id=NULLIF(?,0), updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, rep.Date, rep.VIN, rep.Description, rep.Price, rep.CategoryID, rep.ID)
	return err
}

func (r *RepairRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM repairs WHERE id = ?`, id)
	return err
}
