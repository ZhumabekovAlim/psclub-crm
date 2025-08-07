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

func (r *SettingsRepository) Get(ctx context.Context, companyID, branchID int) (*models.Settings, error) {
	// Получить текущие настройки + имя текущей платежной системы
	query := `
               SELECT s.id, s.payment_type, s.block_time, s.bonus_percent, s.work_time_from, s.work_time_to, s.tables_count, s.notification_time, s.company_id, s.branch_id
               FROM settings s
               WHERE s.company_id=? AND s.branch_id=?
               LIMIT 1
        `
	var s models.Settings
	err := r.db.QueryRowContext(ctx, query, companyID, branchID).Scan(
		&s.ID, &s.PaymentType, &s.BlockTime, &s.BonusPercent, &s.WorkTimeFrom, &s.WorkTimeTo, &s.TablesCount, &s.NotificationTime, &s.CompanyID, &s.BranchID,
	)
	if err != nil {
		return nil, err
	}

	// Получить список всех payment_types
	ptQuery := `SELECT id, name, hold_percent FROM payment_types WHERE company_id=? AND branch_id=? ORDER BY id`
	rows, err := r.db.QueryContext(ctx, ptQuery, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []models.PaymentType
	for rows.Next() {
		var pt models.PaymentType
		if err := rows.Scan(&pt.ID, &pt.Name, &pt.HoldPercent); err != nil {
			return nil, err
		}
		types = append(types, pt)
	}
	s.PaymentTypes = types

	// Получить список всех channels
	chQuery := `SELECT id, company_id, branch_id, name FROM channels WHERE company_id=? AND branch_id=? ORDER BY id`
	chRows, err := r.db.QueryContext(ctx, chQuery, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer chRows.Close()

	var channels []models.Channel
	for chRows.Next() {
		var ch models.Channel
		if err := chRows.Scan(&ch.ID, &ch.CompanyID, &ch.BranchID, &ch.Name); err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}
	s.Channels = channels

	return &s, nil
}

func (r *SettingsRepository) Update(ctx context.Context, s *models.Settings) error {
	query := `
               UPDATE settings
               SET payment_type = ?, block_time = ?, bonus_percent = ?, work_time_from = ?, work_time_to = ?, tables_count = ?, notification_time = ?
               WHERE id = ? AND company_id=? AND branch_id=?
       `
	_, err := r.db.ExecContext(ctx, query, s.PaymentType, s.BlockTime, s.BonusPercent, s.WorkTimeFrom, s.WorkTimeTo, s.TablesCount, s.NotificationTime, s.ID, s.CompanyID, s.BranchID)
	return err
}

func (r *SettingsRepository) Create(ctx context.Context, s *models.Settings) (int, error) {
	query := `
               INSERT INTO settings (payment_type, block_time, bonus_percent, work_time_from, work_time_to, tables_count, notification_time, company_id, branch_id)
               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
       `
	res, err := r.db.ExecContext(ctx, query, s.PaymentType, s.BlockTime, s.BonusPercent, s.WorkTimeFrom, s.WorkTimeTo, s.TablesCount, s.NotificationTime, s.CompanyID, s.BranchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *SettingsRepository) Delete(ctx context.Context, id, companyID, branchID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM settings WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}

func (r *SettingsRepository) GetTablesCount(ctx context.Context, companyID, branchID int) (int, error) {
	var cnt int
	err := r.db.QueryRowContext(ctx, `SELECT tables_count FROM settings WHERE company_id=? AND branch_id=? LIMIT 1`, companyID, branchID).Scan(&cnt)
	return cnt, err
}

func (r *SettingsRepository) GetNotificationTime(ctx context.Context, companyID, branchID int) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT notification_time FROM settings WHERE company_id=? AND branch_id=? LIMIT 1`, companyID, branchID).Scan(&n)
	return n, err
}
