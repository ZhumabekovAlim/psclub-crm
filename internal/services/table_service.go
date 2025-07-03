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

func (s *TableService) GetAllTables(ctx context.Context) ([]models.Table, error) {
	return s.repo.GetAll(ctx)
}

func (s *TableService) GetTableByID(id int) (*models.Table, error) {
	return s.repo.GetByID(id)
}

func (s *TableService) UpdateTable(id int, table *models.Table) error {
	return s.repo.Update(id, table)
}

func (s *TableService) DeleteTable(id int) error {
	return s.repo.Delete(id)
}

func (s *TableService) ReorderTable(ctx context.Context, id int, number int) error {
	return s.repo.Reorder(ctx, id, number)
}
