package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type TableCategoryService struct {
	repo *repositories.TableCategoryRepository
}

func NewTableCategoryService(r *repositories.TableCategoryRepository) *TableCategoryService {
	return &TableCategoryService{repo: r}
}

func (s *TableCategoryService) CreateCategory(ctx context.Context, c *models.TableCategory) (int, error) {
	return s.repo.Create(ctx, c)
}

func (s *TableCategoryService) GetAllCategories(ctx context.Context, companyID, branchID int) ([]models.TableCategory, error) {
	return s.repo.GetAll(ctx, companyID, branchID)
}

func (s *TableCategoryService) GetCategoryByID(ctx context.Context, id, companyID, branchID int) (*models.TableCategory, error) {
	return s.repo.GetByID(ctx, id, companyID, branchID)
}

func (s *TableCategoryService) UpdateCategory(ctx context.Context, category *models.TableCategory) error {
	return s.repo.Update(ctx, category)
}

func (s *TableCategoryService) DeleteCategory(ctx context.Context, id, companyID, branchID int) error {
	return s.repo.Delete(ctx, id, companyID, branchID)
}
