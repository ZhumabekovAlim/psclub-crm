package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

// PaymentTypeRepository handles CRUD operations for payment types.
type PaymentTypeRepository struct {
	db *sql.DB
}

func NewPaymentTypeRepository(db *sql.DB) *PaymentTypeRepository {
	return &PaymentTypeRepository{db: db}
}

func (r *PaymentTypeRepository) GetAll(ctx context.Context) ([]models.PaymentType, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM payment_types ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.PaymentType
	for rows.Next() {
		var pt models.PaymentType
		if err := rows.Scan(&pt.ID, &pt.Name); err != nil {
			return nil, err
		}
		result = append(result, pt)
	}
	return result, nil
}

func (r *PaymentTypeRepository) Create(ctx context.Context, pt *models.PaymentType) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO payment_types (name) VALUES (?)`, pt.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *PaymentTypeRepository) Update(ctx context.Context, pt *models.PaymentType) error {
	_, err := r.db.ExecContext(ctx, `UPDATE payment_types SET name=? WHERE id=?`, pt.Name, pt.ID)
	return err
}

func (r *PaymentTypeRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM payment_types WHERE id=?`, id)
	return err
}
