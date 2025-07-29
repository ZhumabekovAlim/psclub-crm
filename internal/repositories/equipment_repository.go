package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type EquipmentRepository struct {
	db *sql.DB
}

func NewEquipmentRepository(db *sql.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) Create(ctx context.Context, e *models.Equipment) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO equipment (name, quantity, description) VALUES (?, ?, ?)`, e.Name, e.Quantity, e.Description)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *EquipmentRepository) GetAll(ctx context.Context) ([]models.Equipment, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, quantity, IFNULL(description,'') FROM equipment ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Equipment
	for rows.Next() {
		var e models.Equipment
		if err := rows.Scan(&e.ID, &e.Name, &e.Quantity, &e.Description); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *EquipmentRepository) GetByID(ctx context.Context, id int) (*models.Equipment, error) {
	var e models.Equipment
	err := r.db.QueryRowContext(ctx, `SELECT id, name, quantity, IFNULL(description,'') FROM equipment WHERE id=?`, id).Scan(&e.ID, &e.Name, &e.Quantity, &e.Description)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, e *models.Equipment) error {
	_, err := r.db.ExecContext(ctx, `UPDATE equipment SET name=?,  description=? WHERE id=?`, e.Name, e.Description, e.ID)
	return err
}

func (r *EquipmentRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM equipment WHERE id=?`, id)
	return err
}

func (r *EquipmentRepository) SetQuantity(ctx context.Context, id int, qty float64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE equipment SET quantity=? WHERE id=?`, qty, id)
	return err
}
