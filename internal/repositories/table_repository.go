package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) Create(ctx context.Context, t *models.Table) (int, error) {
	// Determine the next table number automatically within company and branch
	err := r.db.QueryRowContext(ctx, `SELECT IFNULL(MAX(number), 0) + 1 FROM tables WHERE company_id=? AND branch_id=?`, t.CompanyID, t.BranchID).Scan(&t.Number)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO tables (category_id, name, number, company_id, branch_id) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, t.CategoryID, t.Name, t.Number, t.CompanyID, t.BranchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *TableRepository) GetAll(ctx context.Context, companyID, branchID int) ([]models.Table, error) {
	query := `SELECT id, category_id, name, number, company_id, branch_id FROM tables WHERE company_id=? AND branch_id=? ORDER BY number`
	rows, err := r.db.QueryContext(ctx, query, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Table
	for rows.Next() {
		var t models.Table
		err := rows.Scan(&t.ID, &t.CategoryID, &t.Name, &t.Number, &t.CompanyID, &t.BranchID)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *TableRepository) GetByID(ctx context.Context, id, companyID, branchID int) (*models.Table, error) {
	var table models.Table
	err := r.db.QueryRowContext(ctx, "SELECT id, name, category_id, number, company_id, branch_id FROM tables WHERE id = ? AND company_id=? AND branch_id=?", id, companyID, branchID).
		Scan(&table.ID, &table.Name, &table.CategoryID, &table.Number, &table.CompanyID, &table.BranchID)
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *TableRepository) Update(ctx context.Context, table *models.Table) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tables SET name = ?, category_id = ?, number = ? WHERE id = ? AND company_id=? AND branch_id=?", table.Name, table.CategoryID, table.Number, table.ID, table.CompanyID, table.BranchID)
	return err
}

func (r *TableRepository) Delete(ctx context.Context, id, companyID, branchID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tables WHERE id = ? AND company_id=? AND branch_id=?", id, companyID, branchID)
	return err
}

func (r *TableRepository) Reorder(ctx context.Context, id, newNumber, companyID, branchID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var current int
	if err := tx.QueryRowContext(ctx, `SELECT number FROM tables WHERE id = ? AND company_id=? AND branch_id=?`, id, companyID, branchID).Scan(&current); err != nil {
		tx.Rollback()
		return err
	}
	if newNumber == current {
		return tx.Commit()
	}
	if newNumber < current {
		if _, err := tx.ExecContext(ctx, `UPDATE tables SET number = number + 1 WHERE number >= ? AND number < ? AND company_id=? AND branch_id=?`, newNumber, current, companyID, branchID); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if _, err := tx.ExecContext(ctx, `UPDATE tables SET number = number - 1 WHERE number > ? AND number <= ? AND company_id=? AND branch_id=?`, current, newNumber, companyID, branchID); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `UPDATE tables SET number = ? WHERE id = ? AND company_id=? AND branch_id=?`, newNumber, id, companyID, branchID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
