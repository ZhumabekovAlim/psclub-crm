package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type EquipmentService struct {
	repo *repositories.EquipmentRepository
}

func NewEquipmentService(r *repositories.EquipmentRepository) *EquipmentService {
	return &EquipmentService{repo: r}
}

func (s *EquipmentService) Create(ctx context.Context, e *models.Equipment) (int, error) {
	return s.repo.Create(ctx, e)
}

func (s *EquipmentService) GetAll(ctx context.Context) ([]models.Equipment, error) {
	return s.repo.GetAll(ctx)
}

func (s *EquipmentService) GetByID(ctx context.Context, id int) (*models.Equipment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EquipmentService) Update(ctx context.Context, e *models.Equipment) error {
	return s.repo.Update(ctx, e)
}

func (s *EquipmentService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *EquipmentService) SetQuantity(ctx context.Context, id int, qty float64) error {
	return s.repo.SetQuantity(ctx, id, qty)
}
