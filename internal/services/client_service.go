package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type ClientService struct {
	repo *repositories.ClientRepository
}

func NewClientService(r *repositories.ClientRepository) *ClientService {
	return &ClientService{repo: r}
}

func (s *ClientService) CreateClient(ctx context.Context, client *models.Client) (int, error) {
	return s.repo.Create(ctx, client)
}

func (s *ClientService) GetAllClients(ctx context.Context) ([]models.Client, error) {
	return s.repo.GetAll(ctx)
}

func (s *ClientService) GetClientByID(ctx context.Context, id int) (*models.Client, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ClientService) UpdateClient(ctx context.Context, client *models.Client) error {
	return s.repo.Update(ctx, client)
}

func (s *ClientService) DeleteClient(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
