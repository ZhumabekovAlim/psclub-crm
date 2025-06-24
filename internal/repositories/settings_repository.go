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
                SELECT s.id, s.payment_type, s.block_time, s.bonus_percent, s.work_time_from, s.work_time_to, s.tables_count, s.notification_time
                FROM settings s
                LEFT JOIN payment_types pt ON s.payment_type = pt.id
                LIMIT 1
        `
	var s models.Settings
	err := r.db.QueryRowContext(ctx, query).Scan(
		&s.ID, &s.PaymentType, &s.BlockTime, &s.BonusPercent, &s.WorkTimeFrom, &s.WorkTimeTo, &s.TablesCount, &s.NotificationTime,
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
                SET payment_type = ?, block_time = ?, bonus_percent = ?, work_time_from = ?, work_time_to = ?, tables_count = ?, notification_time = ?
                WHERE id = ?
        `
	_, err := r.db.ExecContext(ctx, query, s.PaymentType, s.BlockTime, s.BonusPercent, s.WorkTimeFrom, s.WorkTimeTo, s.TablesCount, s.NotificationTime, s.ID)
	return err
}

func (r *SettingsRepository) Create(ctx context.Context, s *models.Settings) (int, error) {
	query := `
                INSERT INTO settings (payment_type, block_time, bonus_percent, work_time_from, work_time_to, tables_count, notification_time)
                VALUES (?, ?, ?, ?, ?, ?, ?)
        `
	res, err := r.db.ExecContext(ctx, query, s.PaymentType, s.BlockTime, s.BonusPercent, s.WorkTimeFrom, s.WorkTimeTo, s.TablesCount, s.NotificationTime)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *SettingsRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM settings WHERE id=?`, id)
	return err
}

func (r *SettingsRepository) GetTablesCount(ctx context.Context) (int, error) {
	var cnt int
	err := r.db.QueryRowContext(ctx, `SELECT tables_count FROM settings LIMIT 1`).Scan(&cnt)
	return cnt, err
}

func (r *SettingsRepository) GetNotificationTime(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT notification_time FROM settings LIMIT 1`).Scan(&n)
	return n, err
}
