package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(r *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) CreateCategory(ctx context.Context, c *models.Category) (int, error) {
	return s.repo.Create(ctx, c)
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, c *models.Category) error {
	return s.repo.Update(ctx, c)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
