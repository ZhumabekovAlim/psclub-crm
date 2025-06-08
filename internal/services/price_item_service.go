package services

import (
	"context"
	"errors"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type PriceItemService struct {
	repo        *repositories.PriceItemRepository
	historyRepo *repositories.PriceItemHistoryRepository
}

func NewPriceItemService(r *repositories.PriceItemRepository, hr *repositories.PriceItemHistoryRepository) *PriceItemService {
	return &PriceItemService{repo: r, historyRepo: hr}
}

func (s *PriceItemService) CreatePriceItem(ctx context.Context, item *models.PriceItem) (int, error) {
	return s.repo.Create(ctx, item)
}

func (s *PriceItemService) GetAllPriceItems(ctx context.Context) ([]models.PriceItem, error) {
	return s.repo.GetAll(ctx)
}

func (s *PriceItemService) GetPriceItemByID(ctx context.Context, id int) (*models.PriceItem, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PriceItemService) UpdatePriceItem(ctx context.Context, item *models.PriceItem) error {
	return s.repo.Update(ctx, item)
}

func (s *PriceItemService) DeletePriceItem(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// Приход товара на склад — создаёт запись в истории и увеличивает остаток
func (s *PriceItemService) AddIncome(ctx context.Context, history *models.PriceItemHistory) error {
	if history.Operation != "INCOME" {
		return errors.New("operation must be INCOME")
	}
	// 1. Добавляем запись в историю
	_, err := s.historyRepo.Create(ctx, history)
	if err != nil {
		return err
	}
	// 2. Увеличиваем остаток в PriceItem
	return s.repo.IncreaseStock(ctx, history.PriceItemID, history.Quantity)
}

// Списание/Продажа товара — запись в истории и уменьшение остатка
func (s *PriceItemService) AddOutcome(ctx context.Context, history *models.PriceItemHistory) error {
	if history.Operation != "OUTCOME" {
		return errors.New("operation must be OUTCOME")
	}
	_, err := s.historyRepo.Create(ctx, history)
	if err != nil {
		return err
	}
	return s.repo.DecreaseStock(ctx, history.PriceItemID, history.Quantity)
}

// Получить историю по товару
func (s *PriceItemService) GetHistoryByItem(ctx context.Context, priceItemID int) ([]models.PriceItemHistory, error) {
	return s.historyRepo.GetByItem(ctx, priceItemID)
}

// Получить всю историю операций
func (s *PriceItemService) GetAllHistory(ctx context.Context) ([]models.PriceItemHistory, error) {
	return s.historyRepo.GetAll(ctx)
}
