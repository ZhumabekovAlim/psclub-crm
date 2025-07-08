package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type EquipmentInventoryHistoryRepository struct {
	db *sql.DB
}

func NewEquipmentInventoryHistoryRepository(db *sql.DB) *EquipmentInventoryHistoryRepository {
	return &EquipmentInventoryHistoryRepository{db: db}
}

func (r *EquipmentInventoryHistoryRepository) Create(ctx context.Context, h *models.EquipmentInventoryHistory) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO equipment_inventory_history (equipment_id, expected, actual, difference, created_at) VALUES (?, ?, ?, ?, NOW())`, h.EquipmentID, h.Expected, h.Actual, h.Difference)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *EquipmentInventoryHistoryRepository) GetAll(ctx context.Context) ([]models.EquipmentInventoryHistory, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT eih.id, eih.equipment_id, e.name, eih.expected, eih.actual, eih.difference, eih.created_at
        FROM equipment_inventory_history eih
        LEFT JOIN equipment e ON eih.equipment_id = e.id
        ORDER BY eih.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.EquipmentInventoryHistory
	for rows.Next() {
		var h models.EquipmentInventoryHistory
		if err := rows.Scan(&h.ID, &h.EquipmentID, &h.Name, &h.Expected, &h.Actual, &h.Difference, &h.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, h)
	}
	return list, nil
}
