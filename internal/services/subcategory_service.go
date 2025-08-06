package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type SubcategoryService struct {
	repo *repositories.SubcategoryRepository
}

func NewSubcategoryService(r *repositories.SubcategoryRepository) *SubcategoryService {
	return &SubcategoryService{repo: r}
}

func (s *SubcategoryService) CreateSubcategory(ctx context.Context, sub *models.Subcategory) (int, error) {
	return s.repo.Create(ctx, sub)
}

func (s *SubcategoryService) GetAllSubcategories(ctx context.Context, companyID, branchID int) ([]models.Subcategory, error) {
	return s.repo.GetAll(ctx, companyID, branchID)
}

func (s *SubcategoryService) GetSubcategoryByID(ctx context.Context, id, companyID, branchID int) (*models.Subcategory, error) {
	return s.repo.GetByID(ctx, id, companyID, branchID)
}

func (s *SubcategoryService) GetSubcategoryByCategoryID(ctx context.Context, id, companyID, branchID int) ([]models.Subcategory, error) {
	return s.repo.GetSubcategoriesByCategoryID(ctx, id, companyID, branchID)
}

func (s *SubcategoryService) UpdateSubcategory(ctx context.Context, sub *models.Subcategory) error {
	return s.repo.Update(ctx, sub)
}

func (s *SubcategoryService) DeleteSubcategory(ctx context.Context, id, companyID, branchID int) error {
	return s.repo.Delete(ctx, id, companyID, branchID)
}
