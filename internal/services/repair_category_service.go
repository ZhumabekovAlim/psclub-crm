package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type RepairCategoryService struct {
	repo *repositories.RepairCategoryRepository
}

func NewRepairCategoryService(r *repositories.RepairCategoryRepository) *RepairCategoryService {
	return &RepairCategoryService{repo: r}
}

func (s *RepairCategoryService) Create(ctx context.Context, c *models.RepairCategory) (int, error) {
	return s.repo.Create(ctx, c)
}

func (s *RepairCategoryService) GetAll(ctx context.Context) ([]models.RepairCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *RepairCategoryService) Update(ctx context.Context, c *models.RepairCategory) error {
	return s.repo.Update(ctx, c)
}

func (s *RepairCategoryService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *RepairCategoryService) GetByName(ctx context.Context, name string) (*models.RepairCategory, error) {
	return s.repo.GetByName(ctx, name)
}
