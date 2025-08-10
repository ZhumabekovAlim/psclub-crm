package services

import (
	"context"
	"errors"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type PriceItemService struct {
	repo          *repositories.PriceItemRepository
	historyRepo   *repositories.PriceItemHistoryRepository
	plHistoryRepo *repositories.PricelistHistoryRepository
}

func NewPriceItemService(r *repositories.PriceItemRepository, hr *repositories.PriceItemHistoryRepository, plhr *repositories.PricelistHistoryRepository) *PriceItemService {
	return &PriceItemService{repo: r, historyRepo: hr, plHistoryRepo: plhr}
}

func (s *PriceItemService) CreatePriceItem(ctx context.Context, item *models.PriceItem) (int, error) {
	ex, err := s.repo.GetByName(ctx, item.Name)
	if err != nil {
		return 0, err
	}
	if ex != nil {
		return 0, ErrNameExists
	}
	return s.repo.Create(ctx, item)
}

func (s *PriceItemService) GetAllPriceItems(ctx context.Context) ([]models.PriceItem, error) {
	return s.repo.GetAll(ctx)
}

func (s *PriceItemService) GetPriceItemByID(ctx context.Context, id int) (*models.PriceItem, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PriceItemService) UpdatePriceItem(ctx context.Context, item *models.PriceItem) error {
	ex, err := s.repo.GetByName(ctx, item.Name)
	if err != nil {
		return err
	}
	if ex != nil && ex.ID != item.ID {
		return ErrNameExists
	}
	return s.repo.Update(ctx, item)
}

func (s *PriceItemService) DeletePriceItem(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *PriceItemService) GetPriceItemsByCategoryName(ctx context.Context, categoryName string) ([]models.PriceItem, error) {
	return s.repo.GetByCategoryName(ctx, categoryName)
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

// CreatePricelistHistory saves a replenishment record without changing stock
func (s *PriceItemService) CreatePricelistHistory(ctx context.Context, hist *models.PricelistHistory) (int, error) {
	if hist == nil {
		return 0, errors.New("history is nil")
	}
	return s.plHistoryRepo.Create(ctx, hist)
}

// Replenish increases stock and saves record to pricelist_history
func (s *PriceItemService) Replenish(ctx context.Context, hist *models.PricelistHistory) error {
	if _, err := s.plHistoryRepo.Create(ctx, hist); err != nil {
		return err
	}

	if err := s.repo.UpdateBuyPrice(ctx, hist.PriceItemID, hist.BuyPrice); err != nil {
		return err
	}

	return s.repo.IncreaseStock(ctx, hist.PriceItemID, float64(hist.Quantity))
}

// GetPricelistHistoryByItem returns replenish history for one price item
func (s *PriceItemService) GetPricelistHistoryByItem(ctx context.Context, id int) ([]models.PricelistHistory, error) {
	return s.plHistoryRepo.GetByItem(ctx, id)
}

func (s *PriceItemService) GetPricelistHistoryByCategory(ctx context.Context, categoryID int) ([]models.PricelistHistory, error) {
	return s.plHistoryRepo.GetByCategory(ctx, categoryID)
}

// GetAllPricelistHistory returns all replenish records
func (s *PriceItemService) GetAllPricelistHistory(ctx context.Context) ([]models.PricelistHistory, error) {
	return s.plHistoryRepo.GetAll(ctx)
}

func (s *PriceItemService) GetPricelistHistoryByID(ctx context.Context, id int) (*models.PricelistHistory, error) {
	return s.plHistoryRepo.GetByID(ctx, id)
}

func (s *PriceItemService) DeletePricelistHistory(ctx context.Context, id int) error {
	return s.plHistoryRepo.Delete(ctx, id)
}
