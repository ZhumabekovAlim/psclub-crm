package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Get(ctx context.Context) (*models.Settings, error) {
	// Получить текущие настройки + имя текущей платежной системы
	query := `
		SELECT s.id, s.payment_type, s.block_time, s.bonus_percent, s.work_time_from, s.work_time_to
		FROM settings s 	
		JOIN payment_types pt ON s.payment_type = pt.id
		LIMIT 1
	`
	var s models.Settings
	err := r.db.QueryRowContext(ctx, query).Scan(
		&s.ID, &s.PaymentType, &s.BlockTime, &s.BonusPercent, &s.WorkTimeFrom, &s.WorkTimeTo,
	)
	if err != nil {
		return nil, err
	}

	// Получить список всех payment_types
	ptQuery := `SELECT id, name FROM payment_types ORDER BY id`
	rows, err := r.db.QueryContext(ctx, ptQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.PaymentType
	for rows.Next() {
		var pt models.PaymentType
		if err := rows.Scan(&pt.ID, &pt.Name); err != nil {
			return nil, err
		}
		types = append(types, pt)
	}
	s.PaymentTypes = types

	return &s, nil
}

func (r *SettingsRepository) Update(ctx context.Context, s *models.Settings) error {
	query := `
		UPDATE settings 
		SET payment_type = ?, block_time = ?, bonus_percent = ?, work_time_from = ?, work_time_to = ? 
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, s.PaymentType, s.BlockTime, s.BonusPercent, s.WorkTimeFrom, s.WorkTimeTo, s.ID)
	return err
}
