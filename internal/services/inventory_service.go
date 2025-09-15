package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type InventoryItem struct {
	PriceItemID int     `json:"price_item_id"`
	Actual      float64 `json:"actual"`
}

type InventoryService struct {
	priceRepo   *repositories.PriceItemRepository
	historyRepo *repositories.InventoryHistoryRepository
	expenseSvc  *ExpenseService
	expCatSvc   *ExpenseCategoryService
}

func NewInventoryService(pr *repositories.PriceItemRepository, hr *repositories.InventoryHistoryRepository, es *ExpenseService, ec *ExpenseCategoryService) *InventoryService {
	return &InventoryService{priceRepo: pr, historyRepo: hr, expenseSvc: es, expCatSvc: ec}
}

func (s *InventoryService) PerformInventory(ctx context.Context, items []InventoryItem) error {
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	var catID int
	if cat, _ := s.expCatSvc.GetByName(ctx, "Инвентаризация"); cat != nil {
		catID = cat.ID
	} else {
		newCat := models.ExpenseCategory{Name: "Инвентаризация", CompanyID: companyID, BranchID: branchID}
		catID, _ = s.expCatSvc.Create(ctx, &newCat)
	}
	for _, it := range items {
		pi, err := s.priceRepo.GetByID(ctx, it.PriceItemID)
		if err != nil {
			return err
		}
		diff := it.Actual - pi.Quantity
		hist := models.InventoryHistory{
			PriceItemID: it.PriceItemID,
			Expected:    pi.Quantity,
			Actual:      it.Actual,
			Difference:  diff,
			CreatedAt:   time.Now(),
		}
		if _, err := s.historyRepo.Create(ctx, &hist); err != nil {
			return err
		}
		if diff < 0 {
			// shortage -> record expense
			exp := models.Expense{
				Date:        time.Now(),
				Title:       "Инвентаризация: " + pi.Name,
				Total:       -diff * pi.BuyPrice,
				Description: "Недостача товара " + pi.Name + " (" + fmt.Sprintf("%.0f", math.Abs(diff)) + " шт.)",
				Paid:        false,
				CategoryID:  catID,
			}
			if _, err := s.expenseSvc.CreateExpense(ctx, &exp); err != nil {
				return err
			}
			// adjust stock to actual
			if err := s.priceRepo.SetStock(ctx, it.PriceItemID, it.Actual); err != nil {
				return err
			}
		} else if diff > 0 {
			// excess - increase stock
			if err := s.priceRepo.SetStock(ctx, it.PriceItemID, it.Actual); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *InventoryService) GetHistory(ctx context.Context, companyID, branchID int) ([]models.InventoryHistory, error) {
	return s.historyRepo.GetAll(ctx, companyID, branchID)
}
