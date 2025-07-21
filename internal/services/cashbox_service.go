package services

import (
	"context"
	"time"

	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

// CashboxService handles operations with cashbox and its history
// Inventory operation sets amount to zero and saves record
// Replenish operation increases amount and records an expense

type CashboxService struct {
	repo       *repositories.CashboxRepository
	histRepo   *repositories.CashboxHistoryRepository
	expenseSvc *ExpenseService
	expCatSvc  *ExpenseCategoryService
}

func NewCashboxService(r *repositories.CashboxRepository, hr *repositories.CashboxHistoryRepository, es *ExpenseService, ec *ExpenseCategoryService) *CashboxService {
	return &CashboxService{repo: r, histRepo: hr, expenseSvc: es, expCatSvc: ec}
}

func (s *CashboxService) GetCashbox(ctx context.Context) (*models.Cashbox, error) {
	return s.repo.Get(ctx)
}

func (s *CashboxService) UpdateCashbox(ctx context.Context, c *models.Cashbox) error {
	return s.repo.Update(ctx, c)
}

// Inventory sets cashbox amount to zero and saves history record
func (s *CashboxService) Inventory(ctx context.Context) error {
	return s.InventoryAmount(ctx, -1)
}

// InventoryAmount performs inventory for a specified amount. If amount is less
// than or equal to zero or greater than current cashbox amount, the whole
// amount is inventoried.
func (s *CashboxService) InventoryAmount(ctx context.Context, amount float64) error {
	box, err := s.repo.Get(ctx)
	if err != nil {
		return err
	}
	if amount <= 0 || amount > box.Amount {
		amount = box.Amount
	}
	hist := models.CashboxHistory{
		Operation: "Инвентаризация",
		Amount:    amount,
	}
	if _, err := s.histRepo.Create(ctx, &hist); err != nil {
		return err
	}
	box.Amount -= amount
	return s.repo.Update(ctx, box)
}

// AddIncome increases cashbox amount and records history without creating an expense entry
func (s *CashboxService) AddIncome(ctx context.Context, amount float64) error {
	box, err := s.repo.Get(ctx)
	if err != nil {
		return err
	}
	box.Amount += amount
	if err := s.repo.Update(ctx, box); err != nil {
		return err
	}
	hist := models.CashboxHistory{
		Operation: "Оплата брони",
		Amount:    amount,
	}
	if _, err := s.histRepo.Create(ctx, &hist); err != nil {
		return err
	}
	return nil
}

// RemoveIncome decreases cashbox amount and records history entry.
func (s *CashboxService) RemoveIncome(ctx context.Context, amount float64) error {
	box, err := s.repo.Get(ctx)
	if err != nil {
		return err
	}
	box.Amount -= amount
	if err := s.repo.Update(ctx, box); err != nil {
		return err
	}
	hist := models.CashboxHistory{
		Operation: "Возврат брони",
		Amount:    -amount,
	}
	if _, err := s.histRepo.Create(ctx, &hist); err != nil {
		return err
	}
	return nil
}

// Replenish adds money to cashbox, records history and creates expense entry
func (s *CashboxService) Replenish(ctx context.Context, amount float64) error {
	box, err := s.repo.Get(ctx)
	if err != nil {
		return err
	}
	box.Amount += amount
	if err := s.repo.Update(ctx, box); err != nil {
		return err
	}
	hist := models.CashboxHistory{
		Operation: "Пополнение",
		Amount:    amount,
	}
	if _, err := s.histRepo.Create(ctx, &hist); err != nil {
		return err
	}

	var catID int
	if cat, _ := s.expCatSvc.GetByName(ctx, "Касса"); cat != nil {
		catID = cat.ID
	} else {
		newCat := models.ExpenseCategory{Name: "Касса"}
		catID, _ = s.expCatSvc.Create(ctx, &newCat)
	}
	exp := models.Expense{
		Date:       time.Now(),
		Title:      "Пополнение кассы",
		Total:      amount,
		Paid:       false,
		CategoryID: catID,
	}
	_, _ = s.expenseSvc.CreateExpense(ctx, &exp)
	return nil
}

func (s *CashboxService) GetHistory(ctx context.Context) ([]models.CashboxHistory, error) {
	return s.histRepo.GetAll(ctx)
}
