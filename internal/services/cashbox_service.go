package services

import (
	"context"
	"fmt"
	"time"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

// CashboxService handles operations with cashbox and its history
// Inventory operation sets amount to zero and saves record
// Replenish operation increases amount and records an expense

type CashboxService struct {
	repo         *repositories.CashboxRepository
	histRepo     *repositories.CashboxHistoryRepository
	expenseSvc   *ExpenseService
	expCatSvc    *ExpenseCategoryService
	settingsRepo *repositories.SettingsRepository
}

func NewCashboxService(r *repositories.CashboxRepository, hr *repositories.CashboxHistoryRepository, es *ExpenseService, ec *ExpenseCategoryService, sr *repositories.SettingsRepository) *CashboxService {
	return &CashboxService{repo: r, histRepo: hr, expenseSvc: es, expCatSvc: ec, settingsRepo: sr}
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

	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)
	var catID int
	if cat, _ := s.expCatSvc.GetByName(ctx, "Касса"); cat != nil {
		catID = cat.ID
	} else {
		newCat := models.ExpenseCategory{Name: "Касса", CompanyID: companyID, BranchID: branchID}
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

// getWorkDayRange calculates current work day start and end based on settings
func getWorkDayRange(now time.Time, fromStr, toStr string) (time.Time, time.Time) {
	loc := now.Location()
	wf, err := time.ParseInLocation("15:04:05", fromStr, loc)
	if err != nil {
		wf = time.Date(0, 1, 1, 0, 0, 0, 0, loc)
	}
	wt, err := time.ParseInLocation("15:04:05", toStr, loc)
	if err != nil {
		wt = time.Date(0, 1, 1, 23, 59, 59, 0, loc)
	}

	start := time.Date(now.Year(), now.Month(), now.Day(), wf.Hour(), wf.Minute(), wf.Second(), 0, loc)
	if now.Before(start) {
		start = start.AddDate(0, 0, -1)
	}

	end := time.Date(start.Year(), start.Month(), start.Day(), wt.Hour(), wt.Minute(), wt.Second(), 0, loc)
	if !end.After(start) {
		end = end.AddDate(0, 0, 1)
	}

	return start, end
}

// GetDay returns cashbox amount at start of current work day and history records for this work day
func (s *CashboxService) GetDay(ctx context.Context) (float64, []models.CashboxHistory, error) {
	now := time.Now()

	fmt.Println("Current time:", now)
	companyID := ctx.Value(common.CtxCompanyID).(int)
	branchID := ctx.Value(common.CtxBranchID).(int)

	settings, err := s.settingsRepo.Get(ctx, companyID, branchID)
	if err != nil {
		return 0, nil, err
	}

	start, end := getWorkDayRange(now, settings.WorkTimeFrom, settings.WorkTimeTo)
	fmt.Println("start: ", start)
	fmt.Println("end: ", end)

	list, err := s.histRepo.GetByPeriod(ctx, start, end)
	if err != nil {
		return 0, nil, err
	}

	box, err := s.repo.Get(ctx)
	if err != nil {
		return 0, nil, err
	}

	startAmount := box.Amount
	for _, h := range list {
		switch h.Operation {
		case "Возврат брони", "Инвентаризация":
			startAmount += h.Amount
		case "Оплата брони", "Пополнение":
			startAmount -= h.Amount
		default:
			startAmount -= h.Amount
		}
	}
	fmt.Println("LIST: ", list)

	return startAmount, list, nil
}
