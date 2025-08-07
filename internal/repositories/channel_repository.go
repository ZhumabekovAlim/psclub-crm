package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
)

type ChannelRepository struct {
	db *sql.DB
}

func NewChannelRepository(db *sql.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Create(ctx context.Context, ch *models.Channel) (int, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	res, err := r.db.ExecContext(ctx, `INSERT INTO channels (name, company_id, branch_id) VALUES (?, ?, ?)`, ch.Name, companyID, branchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ChannelRepository) GetAll(ctx context.Context) ([]models.Channel, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	rows, err := r.db.QueryContext(ctx, `SELECT id, company_id, branch_id, name FROM channels WHERE company_id=? AND branch_id=? ORDER BY id`, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Channel
	for rows.Next() {
		var ch models.Channel
		if err := rows.Scan(&ch.ID, &ch.CompanyID, &ch.BranchID, &ch.Name); err != nil {
			return nil, err
		}
		result = append(result, ch)
	}
	return result, nil
}

func (r *ChannelRepository) Update(ctx context.Context, ch *models.Channel) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `UPDATE channels SET name=? WHERE id=? AND company_id=? AND branch_id=?`, ch.Name, ch.ID, companyID, branchID)
	return err
}

func (r *ChannelRepository) Delete(ctx context.Context, id int) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `DELETE FROM channels WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}

func (r *ChannelRepository) GetByID(ctx context.Context, id int) (*models.Channel, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	var ch models.Channel
	err := r.db.QueryRowContext(ctx, `SELECT id, company_id, branch_id, name FROM channels WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID).Scan(&ch.ID, &ch.CompanyID, &ch.BranchID, &ch.Name)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}
