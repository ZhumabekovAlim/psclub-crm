package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type EquipmentRepository struct {
	db *sql.DB
}

func NewEquipmentRepository(db *sql.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) Create(ctx context.Context, e *models.Equipment) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO equipment (name, quantity, description, company_id, branch_id) VALUES (?, ?, ?, ?, ?)`, e.Name, e.Quantity, e.Description, e.CompanyID, e.BranchID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *EquipmentRepository) GetAll(ctx context.Context, companyID, branchID int) ([]models.Equipment, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, quantity, IFNULL(description,''), company_id, branch_id FROM equipment WHERE company_id=? AND branch_id=? ORDER BY id`, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Equipment
	for rows.Next() {
		var e models.Equipment
		if err := rows.Scan(&e.ID, &e.Name, &e.Quantity, &e.Description, &e.CompanyID, &e.BranchID); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *EquipmentRepository) GetByID(ctx context.Context, id, companyID, branchID int) (*models.Equipment, error) {
	var e models.Equipment
	err := r.db.QueryRowContext(ctx, `SELECT id, name, quantity, IFNULL(description,''), company_id, branch_id FROM equipment WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID).Scan(&e.ID, &e.Name, &e.Quantity, &e.Description, &e.CompanyID, &e.BranchID)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EquipmentRepository) Update(ctx context.Context, e *models.Equipment) error {
	_, err := r.db.ExecContext(ctx, `UPDATE equipment SET name=?,  description=? WHERE id=? AND company_id=? AND branch_id=?`, e.Name, e.Description, e.ID, e.CompanyID, e.BranchID)
	return err
}

func (r *EquipmentRepository) Delete(ctx context.Context, id, companyID, branchID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM equipment WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}

func (r *EquipmentRepository) SetQuantity(ctx context.Context, id int, qty float64, companyID, branchID int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE equipment SET quantity=? WHERE id=? AND company_id=? AND branch_id=?`, qty, id, companyID, branchID)
	return err
}
