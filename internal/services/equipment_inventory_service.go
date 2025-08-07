package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type EquipmentInventoryItem struct {
	EquipmentID int     `json:"equipment_id"`
	Actual      float64 `json:"actual"`
}

type EquipmentInventoryService struct {
	repo       *repositories.EquipmentRepository
	history    *repositories.EquipmentInventoryHistoryRepository
	expenseSvc *ExpenseService
	expCatSvc  *ExpenseCategoryService
}

func NewEquipmentInventoryService(r *repositories.EquipmentRepository, hr *repositories.EquipmentInventoryHistoryRepository, es *ExpenseService, ec *ExpenseCategoryService) *EquipmentInventoryService {
	return &EquipmentInventoryService{repo: r, history: hr, expenseSvc: es, expCatSvc: ec}
}

func (s *EquipmentInventoryService) PerformInventory(ctx context.Context, items []EquipmentInventoryItem, companyID, branchID int) error {
	var catID int
	if cat, _ := s.expCatSvc.GetByName(ctx, "Инвентаризация"); cat != nil {
		catID = cat.ID
	} else {
		newCat := models.ExpenseCategory{Name: "Инвентаризация"}
		catID, _ = s.expCatSvc.Create(ctx, &newCat)
	}
	for _, it := range items {
		eq, err := s.repo.GetByID(ctx, it.EquipmentID, companyID, branchID)
		if err != nil {
			return err
		}
		diff := it.Actual - eq.Quantity
		hist := models.EquipmentInventoryHistory{
			EquipmentID: it.EquipmentID,
			Expected:    eq.Quantity,
			Actual:      it.Actual,
			Difference:  diff,
			CreatedAt:   time.Now(),
			CompanyID:   companyID,
			BranchID:    branchID,
		}
		if _, err := s.history.Create(ctx, &hist); err != nil {
			return err
		}
		if diff < 0 {
			exp := models.Expense{
				Date:        time.Now(),
				Title:       "Инвентаризация: " + eq.Name,
				Total:       0,
				Description: "Недостача оборудования " + eq.Name + " (" + fmt.Sprintf("%.0f", math.Abs(diff)) + " шт.)",
				Paid:        false,
				CategoryID:  catID,
			}
			if _, err := s.expenseSvc.CreateExpense(ctx, &exp); err != nil {
				return err
			}
		}
		if err := s.repo.SetQuantity(ctx, it.EquipmentID, it.Actual, companyID, branchID); err != nil {
			return err
		}
	}
	return nil
}

func (s *EquipmentInventoryService) GetHistory(ctx context.Context, companyID, branchID int) ([]models.EquipmentInventoryHistory, error) {
	return s.history.GetAll(ctx, companyID, branchID)
}
