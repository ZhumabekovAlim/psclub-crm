package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type PriceItemRepository struct {
	db *sql.DB
}

func NewPriceItemRepository(db *sql.DB) *PriceItemRepository {
	return &PriceItemRepository{db: db}
}

func (r *PriceItemRepository) GetByName(ctx context.Context, name string) (*models.PriceItem, error) {
	query := `SELECT id, name, category_id, subcategory_id, quantity, sale_price, buy_price, is_set FROM price_items WHERE name = ?`
	var p models.PriceItem
	err := r.db.QueryRowContext(ctx, query, name).Scan(&p.ID, &p.Name, &p.CategoryID, &p.SubcategoryID, &p.Quantity, &p.SalePrice, &p.BuyPrice, &p.IsSet)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PriceItemRepository) Create(ctx context.Context, p *models.PriceItem) (int, error) {
	query := `INSERT INTO price_items (name, category_id, subcategory_id, quantity, sale_price, buy_price, is_set)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, p.Name, p.CategoryID, p.SubcategoryID, p.Quantity, p.SalePrice, p.BuyPrice, p.IsSet)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *PriceItemRepository) GetAll(ctx context.Context) ([]models.PriceItem, error) {
	query := `SELECT id, name, category_id, subcategory_id, quantity, sale_price, buy_price, is_set FROM price_items ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.PriceItem
	for rows.Next() {
		var p models.PriceItem
		err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.SubcategoryID, &p.Quantity, &p.SalePrice, &p.BuyPrice, &p.IsSet)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *PriceItemRepository) GetByID(ctx context.Context, id int) (*models.PriceItem, error) {
	query := `SELECT id, name, category_id, subcategory_id, quantity, sale_price, buy_price, is_set FROM price_items WHERE id=?`
	var p models.PriceItem
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.CategoryID, &p.SubcategoryID, &p.Quantity, &p.SalePrice, &p.BuyPrice, &p.IsSet)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PriceItemRepository) Update(ctx context.Context, p *models.PriceItem) error {
	query := `UPDATE price_items SET name=?, category_id=?, subcategory_id=?, quantity=?, sale_price=?, buy_price=?, is_set=? WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, p.Name, p.CategoryID, p.SubcategoryID, p.Quantity, p.SalePrice, p.BuyPrice, p.IsSet, p.ID)
	return err
}

func (r *PriceItemRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM price_items WHERE id=?`, id)
	return err
}

// При пополнении склада увеличиваем остаток
func (r *PriceItemRepository) IncreaseStock(ctx context.Context, id int, amount float64) error {
	query := `UPDATE price_items SET quantity = quantity + ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, amount, id)
	return err
}

// UpdateBuyPrice sets a new buy price for the item.
func (r *PriceItemRepository) UpdateBuyPrice(ctx context.Context, id int, price float64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE price_items SET buy_price=? WHERE id=?`, price, id)
	return err
}

// При продаже/списании уменьшаем остаток
func (r *PriceItemRepository) DecreaseStock(ctx context.Context, id int, amount float64) error {
	query := `UPDATE price_items SET quantity = quantity - ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, amount, id)
	return err
}

// SetStock overrides the current quantity with the provided value.
func (r *PriceItemRepository) SetStock(ctx context.Context, id int, quantity float64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE price_items SET quantity=? WHERE id=?`, quantity, id)
	return err
}

func (r *PriceItemRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.PriceItem, error) {

	query := `SELECT pi.id, pi.name, pi.category_id, subcategory_id, quantity, sale_price, buy_price, is_set, s.name AS subcategory_name
                FROM price_items pi
                JOIN subcategories s ON pi.subcategory_id = s.id
                WHERE pi.category_id = ? ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.PriceItem
	for rows.Next() {
		var p models.PriceItem

		if err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.SubcategoryID, &p.Quantity, &p.SalePrice, &p.BuyPrice, &p.IsSet, &p.SubcategoryName); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
