package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type RepairService struct {
	repo *repositories.RepairRepository
}

func NewRepairService(r *repositories.RepairRepository) *RepairService {
	return &RepairService{repo: r}
}

func (s *RepairService) CreateRepair(ctx context.Context, rep *models.Repair) (int, error) {
	return s.repo.Create(ctx, rep)
}

func (s *RepairService) GetAllRepairs(ctx context.Context) ([]models.Repair, error) {
	return s.repo.GetAll(ctx)
}

func (s *RepairService) GetRepairByID(ctx context.Context, id int) (*models.Repair, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RepairService) UpdateRepair(ctx context.Context, rep *models.Repair) error {
	return s.repo.Update(ctx, rep)
}

func (s *RepairService) DeleteRepair(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
