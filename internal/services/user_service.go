package services

import (
	"context"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(r *repositories.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(ctx context.Context, u *models.User) (int, error) {
	return s.repo.Create(ctx, u)
}

func (s *UserService) GetAllUsers(ctx context.Context, companyID, branchID int) ([]models.User, error) {
	return s.repo.GetAll(ctx, companyID, branchID)
}

func (s *UserService) GetUserByID(ctx context.Context, id, companyID, branchID int) (*models.User, error) {
	return s.repo.GetByID(ctx, id, companyID, branchID)
}

func (s *UserService) UpdateUser(ctx context.Context, u *models.User) error {
	return s.repo.Update(ctx, u)
}

func (s *UserService) DeleteUser(ctx context.Context, id, companyID, branchID int) error {
	return s.repo.Delete(ctx, id, companyID, branchID)
}
