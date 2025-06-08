package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type ExpenseService struct {
	repo *repositories.ExpenseRepository
}

func NewExpenseService(r *repositories.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: r}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, e *models.Expense) (int, error) {
	return s.repo.Create(ctx, e)
}

func (s *ExpenseService) GetAllExpenses(ctx context.Context) ([]models.Expense, error) {
	return s.repo.GetAll(ctx)
}

func (s *ExpenseService) GetExpenseByID(ctx context.Context, id int) (*models.Expense, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, e *models.Expense) error {
	return s.repo.Update(ctx, e)
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
