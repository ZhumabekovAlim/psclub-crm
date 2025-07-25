package repositories

import (
	"context"
	"database/sql"
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
	query := `INSERT INTO price_sets (name, category_id, subcategory_id, price) VALUES (?, ?, ?, ?)`
	args := []any{s.Name, s.CategoryID, s.SubcategoryID, s.Price}
	if s.ID > 0 {
		query = `INSERT INTO price_sets (id, name, category_id, subcategory_id, price) VALUES (?, ?, ?, ?, ?)`
		args = []any{s.ID, s.Name, s.CategoryID, s.SubcategoryID, s.Price}
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
	return int(setID), nil
}

func (r *PriceSetRepository) getItems(ctx context.Context, setID int) ([]models.SetItem, error) {
	query := `SELECT si.id, si.price_set_id, si.item_id, si.quantity, pi.name
                FROM set_items si
                JOIN price_items pi ON si.item_id = pi.id
                WHERE si.price_set_id=?`
	rows, err := r.db.QueryContext(ctx, query, setID)
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
	rows, err := r.db.QueryContext(ctx, `SELECT price_sets.id, price_sets.name, price_sets.category_id, subcategory_id, price, subcategories.name FROM price_sets
                                                    JOIN subcategories ON price_sets.subcategory_id = subcategories.id
                                                    ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sets []models.PriceSet
	for rows.Next() {
		var s models.PriceSet
		if err := rows.Scan(&s.ID, &s.Name, &s.CategoryID, &s.SubcategoryID, &s.Price, &s.SubcategoryName); err != nil {
			return nil, err
		}
		s.Items, _ = r.getItems(ctx, s.ID)
		sets = append(sets, s)
	}
	return sets, nil
}

func (r *PriceSetRepository) GetByID(ctx context.Context, id int) (*models.PriceSet, error) {
	var s models.PriceSet
	err := r.db.QueryRowContext(ctx, `SELECT price_sets.id, price_sets.name, price_sets.category_id, subcategory_id, price, subcategories.name FROM price_sets
                                                    JOIN subcategories ON price_sets.subcategory_id = subcategories.id WHERE price_sets.id = ?`, id).Scan(&s.ID, &s.Name, &s.CategoryID, &s.SubcategoryID, &s.Price, &s.SubcategoryName)
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
	_, err = tx.ExecContext(ctx, `UPDATE price_sets SET name=?, category_id=?, subcategory_id=?, price=? WHERE id=?`, s.Name, s.CategoryID, s.SubcategoryID, s.Price, s.ID)
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
	_, err := r.db.ExecContext(ctx, `DELETE FROM price_sets WHERE id=?`, id)
	return err
}

// GetByItem returns all sets that include the given item.
func (r *PriceSetRepository) GetByItem(ctx context.Context, itemID int) ([]models.PriceSet, error) {
	query := `SELECT ps.id, ps.name, ps.category_id, ps.subcategory_id, ps.price
                 FROM price_sets ps
                 JOIN set_items si ON ps.id = si.price_set_id
                 WHERE si.item_id = ?`
	rows, err := r.db.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sets []models.PriceSet
	for rows.Next() {
		var s models.PriceSet
		if err := rows.Scan(&s.ID, &s.Name, &s.CategoryID, &s.SubcategoryID, &s.Price); err != nil {
			return nil, err
		}
		s.Items, _ = r.getItems(ctx, s.ID)
		sets = append(sets, s)
	}
	return sets, nil
}
