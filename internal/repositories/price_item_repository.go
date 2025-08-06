package repositories

import (
	"context"
	"database/sql"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
)

type PriceItemRepository struct {
	db *sql.DB
}

func NewPriceItemRepository(db *sql.DB) *PriceItemRepository {
	return &PriceItemRepository{db: db}
}

func (r *PriceItemRepository) GetByName(ctx context.Context, name string) (*models.PriceItem, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT id, name, category_id, subcategory_id, quantity, sale_price, buy_price, is_set, company_id, branch_id FROM price_items WHERE name = ? AND company_id = ? AND branch_id = ?`
	var p models.PriceItem
	err := r.db.QueryRowContext(ctx, query, name, companyID, branchID).Scan(&p.ID, &p.Name, &p.CategoryID, &p.SubcategoryID, &p.Quantity, &p.SalePrice, &p.BuyPrice, &p.IsSet, &p.CompanyID, &p.BranchID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PriceItemRepository) Create(ctx context.Context, p *models.PriceItem) (int, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `INSERT INTO price_items (name, category_id, subcategory_id, quantity, sale_price, buy_price, is_set, company_id, branch_id)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, p.Name, p.CategoryID, p.SubcategoryID, p.Quantity, p.SalePrice, p.BuyPrice, p.IsSet, companyID, branchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *PriceItemRepository) GetAll(ctx context.Context) ([]models.PriceItem, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT id, name, category_id, subcategory_id, quantity, sale_price, buy_price, is_set, company_id, branch_id FROM price_items WHERE company_id=? AND branch_id=? ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PriceItem
	for rows.Next() {
		var p models.PriceItem
		err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.SubcategoryID, &p.Quantity, &p.SalePrice, &p.BuyPrice, &p.IsSet, &p.CompanyID, &p.BranchID)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *PriceItemRepository) GetByID(ctx context.Context, id int) (*models.PriceItem, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT id, name, category_id, subcategory_id, quantity, sale_price, buy_price, is_set, company_id, branch_id FROM price_items WHERE id=? AND company_id=? AND branch_id=?`
	var p models.PriceItem
	err := r.db.QueryRowContext(ctx, query, id, companyID, branchID).Scan(&p.ID, &p.Name, &p.CategoryID, &p.SubcategoryID, &p.Quantity, &p.SalePrice, &p.BuyPrice, &p.IsSet, &p.CompanyID, &p.BranchID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PriceItemRepository) Update(ctx context.Context, p *models.PriceItem) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `UPDATE price_items SET name=?, category_id=?, subcategory_id=?, quantity=?, sale_price=?, buy_price=?, is_set=? WHERE id=? AND company_id=? AND branch_id=?`
	_, err := r.db.ExecContext(ctx, query, p.Name, p.CategoryID, p.SubcategoryID, p.Quantity, p.SalePrice, p.BuyPrice, p.IsSet, p.ID, companyID, branchID)
	return err
}

func (r *PriceItemRepository) Delete(ctx context.Context, id int) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `DELETE FROM price_items WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}

// При пополнении склада увеличиваем остаток
func (r *PriceItemRepository) IncreaseStock(ctx context.Context, id int, amount float64) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `UPDATE price_items SET quantity = quantity + ? WHERE id = ? AND company_id=? AND branch_id=?`
	_, err := r.db.ExecContext(ctx, query, amount, id, companyID, branchID)
	return err
}

// UpdateBuyPrice sets a new buy price for the item.
func (r *PriceItemRepository) UpdateBuyPrice(ctx context.Context, id int, price float64) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `UPDATE price_items SET buy_price=? WHERE id=? AND company_id=? AND branch_id=?`, price, id, companyID, branchID)
	return err
}

// При продаже/списании уменьшаем остаток
func (r *PriceItemRepository) DecreaseStock(ctx context.Context, id int, amount float64) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `UPDATE price_items SET quantity = quantity - ? WHERE id = ? AND company_id=? AND branch_id=?`
	_, err := r.db.ExecContext(ctx, query, amount, id, companyID, branchID)
	return err
}

// SetStock overrides the current quantity with the provided value.
func (r *PriceItemRepository) SetStock(ctx context.Context, id int, quantity float64) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `UPDATE price_items SET quantity=? WHERE id=? AND company_id=? AND branch_id=?`, quantity, id, companyID, branchID)
	return err
}

func (r *PriceItemRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.PriceItem, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT pi.id, pi.name, pi.category_id, subcategory_id, quantity, sale_price, buy_price, is_set, s.name AS subcategory_name, pi.company_id, pi.branch_id
               FROM price_items pi
               JOIN subcategories s ON pi.subcategory_id = s.id
               WHERE pi.category_id = ? AND pi.company_id = ? AND pi.branch_id = ? ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, categoryID, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.PriceItem
	for rows.Next() {
		var p models.PriceItem

		if err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.SubcategoryID, &p.Quantity, &p.SalePrice, &p.BuyPrice, &p.IsSet, &p.SubcategoryName, &p.CompanyID, &p.BranchID); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
