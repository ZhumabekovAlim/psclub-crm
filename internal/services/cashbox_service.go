package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type CashboxService struct {
	repo *repositories.CashboxRepository
}

func NewCashboxService(r *repositories.CashboxRepository) *CashboxService {
	return &CashboxService{repo: r}
}

func (s *CashboxService) GetCashbox(ctx context.Context) (*models.Cashbox, error) {
	return s.repo.Get(ctx)
}

func (s *CashboxService) UpdateCashbox(ctx context.Context, c *models.Cashbox) error {
	return s.repo.Update(ctx, c)
}
