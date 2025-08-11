package repositories

import (
	"context"
	"database/sql"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
)

type PricelistHistoryRepository struct {
	db *sql.DB
}

func NewPricelistHistoryRepository(db *sql.DB) *PricelistHistoryRepository {
	return &PricelistHistoryRepository{db: db}
}

func (r *PricelistHistoryRepository) Create(ctx context.Context, h *models.PricelistHistory) (int, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `INSERT INTO pricelist_history (price_item_id, quantity, buy_price, total, user_id, company_id, branch_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	res, err := r.db.ExecContext(ctx, query, h.PriceItemID, h.Quantity, h.BuyPrice, h.Total, h.UserID, companyID, branchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *PricelistHistoryRepository) GetByItem(ctx context.Context, priceItemID int) ([]models.PricelistHistory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT id, price_item_id, quantity, buy_price, total, user_id, company_id, branch_id, created_at FROM pricelist_history WHERE price_item_id = ? AND company_id=? AND branch_id=? ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query, priceItemID, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PricelistHistory
	for rows.Next() {
		var h models.PricelistHistory
		if err := rows.Scan(&h.ID, &h.PriceItemID, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CompanyID, &h.BranchID, &h.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (r *PricelistHistoryRepository) GetAll(ctx context.Context) ([]models.PricelistHistory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT ph.id, pi.name, ph.price_item_id, ph.quantity, ph.buy_price, ph.total, ph.user_id, ph.company_id, ph.branch_id, ph.created_at FROM pricelist_history ph
           JOIN price_items pi on ph.price_item_id = pi.id
           WHERE ph.company_id=? AND ph.branch_id=?
           ORDER BY id DESC`
	rows, err := r.db.QueryContext(ctx, query, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PricelistHistory
	for rows.Next() {
		var h models.PricelistHistory
		if err := rows.Scan(&h.ID, &h.ItemName, &h.PriceItemID, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CompanyID, &h.BranchID, &h.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (r *PricelistHistoryRepository) GetByID(ctx context.Context, id int) (*models.PricelistHistory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT ph.id, pi.name, ph.price_item_id, ph.quantity, ph.buy_price, ph.total, ph.user_id, ph.company_id, ph.branch_id, ph.created_at
               FROM pricelist_history ph
               JOIN price_items pi ON ph.price_item_id = pi.id
               WHERE ph.id = ? AND ph.company_id=? AND ph.branch_id=?`
	var h models.PricelistHistory
	err := r.db.QueryRowContext(ctx, query, id, companyID, branchID).Scan(&h.ID, &h.ItemName, &h.PriceItemID, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CompanyID, &h.BranchID, &h.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *PricelistHistoryRepository) Delete(ctx context.Context, id int) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `DELETE FROM pricelist_history WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}
func (r *PricelistHistoryRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.PricelistHistory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT ph.id, pi.name, ph.price_item_id, ph.quantity, ph.buy_price, ph.total, ph.user_id, ph.company_id, ph.branch_id, ph.created_at, u.name AS user_name
               FROM pricelist_history ph
               JOIN price_items pi ON ph.price_item_id = pi.id
               JOIN users u ON ph.user_id = u.id
               WHERE pi.category_id = ? AND ph.company_id=? AND ph.branch_id=? ORDER BY ph.id DESC`
	rows, err := r.db.QueryContext(ctx, query, categoryID, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PricelistHistory
	for rows.Next() {
		var h models.PricelistHistory
		if err := rows.Scan(&h.ID, &h.ItemName, &h.PriceItemID, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CompanyID, &h.BranchID, &h.CreatedAt, &h.UserName); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (r *PricelistHistoryRepository) GetByCategoryName(ctx context.Context, categoryName string) ([]models.PricelistHistory, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT ph.id, pi.name, ph.price_item_id, ph.quantity, ph.buy_price, ph.total, ph.user_id, ph.company_id, ph.branch_id, ph.created_at, u.name AS user_name
              FROM pricelist_history ph
              JOIN price_items pi ON ph.price_item_id = pi.id
              JOIN categories c ON c.id = pi.category_id
              JOIN users u ON ph.user_id = u.id
              WHERE c.name = ? AND ph.company_id=? AND ph.branch_id=? ORDER BY ph.id DESC`
	rows, err := r.db.QueryContext(ctx, query, categoryName, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PricelistHistory
	for rows.Next() {
		var h models.PricelistHistory
		if err := rows.Scan(&h.ID, &h.ItemName, &h.PriceItemID, &h.Quantity, &h.BuyPrice, &h.Total, &h.UserID, &h.CompanyID, &h.BranchID, &h.CreatedAt, &h.UserName); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}
