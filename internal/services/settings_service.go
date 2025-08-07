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

func (s *SettingsService) GetSettings(ctx context.Context, companyID, branchID int) (*models.Settings, error) {
	return s.repo.Get(ctx, companyID, branchID)
}

func (s *SettingsService) UpdateSettings(ctx context.Context, set *models.Settings) error {
	return s.repo.Update(ctx, set)
}

func (s *SettingsService) CreateSettings(ctx context.Context, set *models.Settings) (int, error) {
	return s.repo.Create(ctx, set)
}

func (s *SettingsService) DeleteSettings(ctx context.Context, id, companyID, branchID int) error {
	return s.repo.Delete(ctx, id, companyID, branchID)
}

func (s *SettingsService) GetTablesCount(ctx context.Context, companyID, branchID int) (int, error) {
	return s.repo.GetTablesCount(ctx, companyID, branchID)
}

func (s *SettingsService) GetNotificationTime(ctx context.Context, companyID, branchID int) (int, error) {
	return s.repo.GetNotificationTime(ctx, companyID, branchID)
}
