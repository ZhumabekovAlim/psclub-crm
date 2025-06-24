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

func (s *SettingsService) CreateSettings(ctx context.Context, set *models.Settings) (int, error) {
	return s.repo.Create(ctx, set)
}

func (s *SettingsService) DeleteSettings(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *SettingsService) GetTablesCount(ctx context.Context) (int, error) {
	return s.repo.GetTablesCount(ctx)
}

func (s *SettingsService) GetNotificationTime(ctx context.Context) (int, error) {
	return s.repo.GetNotificationTime(ctx)
}
