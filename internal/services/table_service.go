package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type TableService struct {
	repo *repositories.TableRepository
}

func NewTableService(r *repositories.TableRepository) *TableService {
	return &TableService{repo: r}
}

func (s *TableService) CreateTable(ctx context.Context, t *models.Table) (int, error) {
	return s.repo.Create(ctx, t)
}

func (s *TableService) GetAllTables(ctx context.Context, companyID, branchID int) ([]models.Table, error) {
	return s.repo.GetAll(ctx, companyID, branchID)
}

func (s *TableService) GetTableByID(ctx context.Context, id, companyID, branchID int) (*models.Table, error) {
	return s.repo.GetByID(ctx, id, companyID, branchID)
}

func (s *TableService) UpdateTable(ctx context.Context, table *models.Table) error {
	return s.repo.Update(ctx, table)
}

func (s *TableService) DeleteTable(ctx context.Context, id, companyID, branchID int) error {
	return s.repo.Delete(ctx, id, companyID, branchID)
}

func (s *TableService) ReorderTable(ctx context.Context, id, number, companyID, branchID int) error {
	return s.repo.Reorder(ctx, id, number, companyID, branchID)
}
