package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type SettingsService struct {
	repo *repositories.SettingsRepository
}

func NewSettingsService(r *repositories.SettingsRepository) *SettingsService {
	return &SettingsService{repo: r}
}

func (s *SettingsService) GetSettings(ctx context.Context) (*models.Settings, error) {
	return s.repo.Get(ctx)
}

func (s *SettingsService) UpdateSettings(ctx context.Context, set *models.Settings) error {
	return s.repo.Update(ctx, set)
}
