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
	query := `SELECT id, payment_type, block_time, bonus_percent, work_time_from, work_time_to FROM settings LIMIT 1`
	var s models.Settings
	err := r.db.QueryRowContext(ctx, query).Scan(&s.ID, &s.PaymentType, &s.BlockTime, &s.BonusPercent, &s.WorkTimeFrom, &s.WorkTimeTo)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SettingsRepository) Update(ctx context.Context, s *models.Settings) error {
	query := `UPDATE settings SET payment_type=?, block_time=?, bonus_percent=?, work_time_from=?, work_time_to=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, s.PaymentType, s.BlockTime, s.BonusPercent, s.WorkTimeFrom, s.WorkTimeTo, s.ID)
	return err
}
