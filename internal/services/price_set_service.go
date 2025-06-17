package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type PriceSetService struct {
	repo *repositories.PriceSetRepository
}

func NewPriceSetService(r *repositories.PriceSetRepository) *PriceSetService {
	return &PriceSetService{repo: r}
}

func (s *PriceSetService) CreatePriceSet(ctx context.Context, ps *models.PriceSet) (int, error) {
	return s.repo.Create(ctx, ps)
}

func (s *PriceSetService) GetAllPriceSets(ctx context.Context) ([]models.PriceSet, error) {
	return s.repo.GetAll(ctx)
}

func (s *PriceSetService) GetPriceSetByID(ctx context.Context, id int) (*models.PriceSet, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PriceSetService) UpdatePriceSet(ctx context.Context, ps *models.PriceSet) error {
	return s.repo.Update(ctx, ps)
}

func (s *PriceSetService) DeletePriceSet(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
