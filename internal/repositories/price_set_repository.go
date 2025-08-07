package repositories

import (
	"context"
	"database/sql"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
)

type PriceSetRepository struct {
	db *sql.DB
}

func NewPriceSetRepository(db *sql.DB) *PriceSetRepository {
	return &PriceSetRepository{db: db}
}

func (r *PriceSetRepository) Create(ctx context.Context, s *models.PriceSet) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `INSERT INTO price_sets (name, category_id, subcategory_id, price, company_id, branch_id) VALUES (?, ?, ?, ?, ?, ?)`
	args := []any{s.Name, s.CategoryID, s.SubcategoryID, s.Price, companyID, branchID}
	if s.ID > 0 {
		query = `INSERT INTO price_sets (id, name, category_id, subcategory_id, price, company_id, branch_id) VALUES (?, ?, ?, ?, ?, ?, ?)`
		args = []any{s.ID, s.Name, s.CategoryID, s.SubcategoryID, s.Price, companyID, branchID}
	}
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	setID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if s.ID > 0 {
		setID = int64(s.ID)
	}
	if len(s.Items) > 0 {
		for _, it := range s.Items {
			_, err = tx.ExecContext(ctx, `INSERT INTO set_items (price_set_id, item_id, quantity) VALUES (?, ?, ?)`, setID, it.ItemID, it.Quantity)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	s.CompanyID = companyID
	s.BranchID = branchID
	return int(setID), nil
}

func (r *PriceSetRepository) getItems(ctx context.Context, setID int) ([]models.SetItem, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT si.id, si.price_set_id, si.item_id, si.quantity, pi.name
               FROM set_items si
               JOIN price_items pi ON si.item_id = pi.id
               JOIN price_sets ps ON si.price_set_id = ps.id
               WHERE si.price_set_id=? AND ps.company_id=? AND ps.branch_id=?`
	rows, err := r.db.QueryContext(ctx, query, setID, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.SetItem
	for rows.Next() {
		var it models.SetItem
		if err := rows.Scan(&it.ID, &it.PriceSetID, &it.ItemID, &it.Quantity, &it.ItemName); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func (r *PriceSetRepository) GetAll(ctx context.Context) ([]models.PriceSet, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	rows, err := r.db.QueryContext(ctx, `SELECT price_sets.id, price_sets.name, price_sets.category_id, subcategory_id, price, subcategories.name, price_sets.company_id, price_sets.branch_id FROM price_sets
                                                   JOIN subcategories ON price_sets.subcategory_id = subcategories.id
                                                   WHERE price_sets.company_id=? AND price_sets.branch_id=?
                                                   ORDER BY id`, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sets []models.PriceSet
	for rows.Next() {
		var s models.PriceSet
		if err := rows.Scan(&s.ID, &s.Name, &s.CategoryID, &s.SubcategoryID, &s.Price, &s.SubcategoryName, &s.CompanyID, &s.BranchID); err != nil {
			return nil, err
		}
		s.Items, _ = r.getItems(ctx, s.ID)
		sets = append(sets, s)
	}
	return sets, nil
}

func (r *PriceSetRepository) GetByID(ctx context.Context, id int) (*models.PriceSet, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	var s models.PriceSet
	err := r.db.QueryRowContext(ctx, `SELECT price_sets.id, price_sets.name, price_sets.category_id, subcategory_id, price, subcategories.name, price_sets.company_id, price_sets.branch_id FROM price_sets
                                                     JOIN subcategories ON price_sets.subcategory_id = subcategories.id WHERE price_sets.id = ? AND price_sets.company_id=? AND price_sets.branch_id=?`, id, companyID, branchID).Scan(&s.ID, &s.Name, &s.CategoryID, &s.SubcategoryID, &s.Price, &s.SubcategoryName, &s.CompanyID, &s.BranchID)
	if err != nil {
		return nil, err
	}
	s.Items, _ = r.getItems(ctx, s.ID)
	return &s, nil
}

func (r *PriceSetRepository) Update(ctx context.Context, s *models.PriceSet) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err = tx.ExecContext(ctx, `UPDATE price_sets SET name=?, category_id=?, subcategory_id=?, price=? WHERE id=? AND company_id=? AND branch_id=?`, s.Name, s.CategoryID, s.SubcategoryID, s.Price, s.ID, companyID, branchID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM set_items WHERE price_set_id=?`, s.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, it := range s.Items {
		_, err = tx.ExecContext(ctx, `INSERT INTO set_items (price_set_id, item_id, quantity) VALUES (?, ?, ?)`, s.ID, it.ItemID, it.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *PriceSetRepository) Delete(ctx context.Context, id int) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	_, err := r.db.ExecContext(ctx, `DELETE FROM price_sets WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}

// GetByItem returns all sets that include the given item.
func (r *PriceSetRepository) GetByItem(ctx context.Context, itemID int) ([]models.PriceSet, error) {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	query := `SELECT ps.id, ps.name, ps.category_id, ps.subcategory_id, ps.price, ps.company_id, ps.branch_id
                FROM price_sets ps
                JOIN set_items si ON ps.id = si.price_set_id
                WHERE si.item_id = ? AND ps.company_id=? AND ps.branch_id=?`
	rows, err := r.db.QueryContext(ctx, query, itemID, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sets []models.PriceSet
	for rows.Next() {
		var s models.PriceSet
		if err := rows.Scan(&s.ID, &s.Name, &s.CategoryID, &s.SubcategoryID, &s.Price, &s.CompanyID, &s.BranchID); err != nil {
			return nil, err
		}
		s.Items, _ = r.getItems(ctx, s.ID)
		sets = append(sets, s)
	}
	return sets, nil
}
