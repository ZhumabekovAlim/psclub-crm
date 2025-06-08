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

func (s *TableCategoryService) GetAllCategories(ctx context.Context) ([]models.TableCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *TableCategoryService) GetCategoryByID(id int) (*models.TableCategory, error) {
	return s.repo.GetByID(id)
}

func (s *TableCategoryService) UpdateCategory(id int, category *models.TableCategory) error {
	return s.repo.Update(id, category)
}

func (s *TableCategoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}
