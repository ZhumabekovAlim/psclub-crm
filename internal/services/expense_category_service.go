package services

import (
	"context"

	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type ExpenseCategoryService struct {
	repo *repositories.ExpenseCategoryRepository
}

func NewExpenseCategoryService(r *repositories.ExpenseCategoryRepository) *ExpenseCategoryService {
	return &ExpenseCategoryService{repo: r}
}

func (s *ExpenseCategoryService) Create(ctx context.Context, c *models.ExpenseCategory) (int, error) {
	ex, err := s.repo.GetByName(ctx, c.Name)
	if err != nil {
		return 0, err
	}
	if ex != nil {
		return 0, ErrNameExists
	}
	return s.repo.Create(ctx, c)
}

func (s *ExpenseCategoryService) GetAll(ctx context.Context) ([]models.ExpenseCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *ExpenseCategoryService) Update(ctx context.Context, c *models.ExpenseCategory) error {
	ex, err := s.repo.GetByName(ctx, c.Name)
	if err != nil {
		return err
	}
	if ex != nil && ex.ID != c.ID {
		return ErrNameExists
	}
	return s.repo.Update(ctx, c)
}

func (s *ExpenseCategoryService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *ExpenseCategoryService) GetByName(ctx context.Context, name string) (*models.ExpenseCategory, error) {
	return s.repo.GetByName(ctx, name)
}
