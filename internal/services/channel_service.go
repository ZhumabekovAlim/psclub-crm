package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type ChannelService struct {
	repo *repositories.ChannelRepository
}

func NewChannelService(r *repositories.ChannelRepository) *ChannelService {
	return &ChannelService{repo: r}
}

func (s *ChannelService) Create(ctx context.Context, ch *models.Channel) (int, error) {
	return s.repo.Create(ctx, ch)
}

func (s *ChannelService) GetAll(ctx context.Context) ([]models.Channel, error) {
	return s.repo.GetAll(ctx)
}

func (s *ChannelService) Update(ctx context.Context, ch *models.Channel) error {
	return s.repo.Update(ctx, ch)
}

func (s *ChannelService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *ChannelService) GetByID(ctx context.Context, id int) (*models.Channel, error) {
	return s.repo.GetByID(ctx, id)
}
