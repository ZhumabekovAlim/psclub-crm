package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

// PaymentTypeService provides business logic for payment types.
type PaymentTypeService struct {
	repo *repositories.PaymentTypeRepository
}

func NewPaymentTypeService(r *repositories.PaymentTypeRepository) *PaymentTypeService {
	return &PaymentTypeService{repo: r}
}

func (s *PaymentTypeService) GetAllPaymentTypes(ctx context.Context) ([]models.PaymentType, error) {
	return s.repo.GetAll(ctx)
}

func (s *PaymentTypeService) CreatePaymentType(ctx context.Context, pt *models.PaymentType) (int, error) {
	return s.repo.Create(ctx, pt)
}

func (s *PaymentTypeService) UpdatePaymentType(ctx context.Context, pt *models.PaymentType) error {
	return s.repo.Update(ctx, pt)
}

func (s *PaymentTypeService) DeletePaymentType(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *PaymentTypeService) GetPaymentTypeByID(ctx context.Context, id int) (*models.PaymentType, error) {
	return s.repo.GetByID(ctx, id)
}
