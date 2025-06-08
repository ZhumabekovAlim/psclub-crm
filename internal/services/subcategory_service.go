package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type SubcategoryService struct {
	repo *repositories.SubcategoryRepository
}

func NewSubcategoryService(r *repositories.SubcategoryRepository) *SubcategoryService {
	return &SubcategoryService{repo: r}
}

func (s *SubcategoryService) CreateSubcategory(ctx context.Context, sub *models.Subcategory) (int, error) {
	return s.repo.Create(ctx, sub)
}

func (s *SubcategoryService) GetAllSubcategories(ctx context.Context) ([]models.Subcategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *SubcategoryService) GetSubcategoryByID(ctx context.Context, id int) (*models.Subcategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubcategoryService) UpdateSubcategory(ctx context.Context, sub *models.Subcategory) error {
	return s.repo.Update(ctx, sub)
}

func (s *SubcategoryService) DeleteSubcategory(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
