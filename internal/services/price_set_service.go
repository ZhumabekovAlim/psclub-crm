package services

import (
	"context"
	"math"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type PriceSetService struct {
	repo         *repositories.PriceSetRepository
	itemRepo     *repositories.PriceItemRepository
	categoryRepo *repositories.CategoryRepository
}

const hoursCategoryName = "\u0427\u0430\u0441\u044b"

func (s *PriceSetService) isHoursCategory(ctx context.Context, categoryID int) (bool, error) {
	cat, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return false, err
	}
	return cat.Name == hoursCategoryName, nil
}

func NewPriceSetService(r *repositories.PriceSetRepository, ir *repositories.PriceItemRepository, cr *repositories.CategoryRepository) *PriceSetService {
	return &PriceSetService{repo: r, itemRepo: ir, categoryRepo: cr}
}

func (s *PriceSetService) CreatePriceSet(ctx context.Context, ps *models.PriceSet) (int, error) {
	item := models.PriceItem{
		Name:          ps.Name,
		CategoryID:    ps.CategoryID,
		SubcategoryID: ps.SubcategoryID,
		SalePrice:     float64(ps.Price),
		IsSet:         true,
	}
	id, err := s.itemRepo.Create(ctx, &item)
	if err != nil {
		return 0, err
	}
	ps.ID = id
	if _, err = s.repo.Create(ctx, ps); err != nil {
		return 0, err
	}
	qty, err := s.calculateQuantity(ctx, ps)
	if err == nil {
		ps.Quantity = float64(qty)
		if err = s.itemRepo.SetStock(ctx, ps.ID, float64(qty)); err != nil {
			return id, err
		}
	}
	return id, nil
}

func (s *PriceSetService) GetAllPriceSets(ctx context.Context) ([]models.PriceSet, error) {
	sets, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	for i := range sets {
		qty, err := s.calculateQuantity(ctx, &sets[i])
		if err != nil {
			return nil, err
		}
		sets[i].Quantity = float64(qty)
	}
	return sets, nil
}

func (s *PriceSetService) GetPriceSetByID(ctx context.Context, id int) (*models.PriceSet, error) {
	set, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	qty, err := s.calculateQuantity(ctx, set)
	if err == nil {
		set.Quantity = float64(qty)
	}
	return set, nil
}

func (s *PriceSetService) UpdatePriceSet(ctx context.Context, ps *models.PriceSet) error {
	item := models.PriceItem{
		ID:            ps.ID,
		Name:          ps.Name,
		CategoryID:    ps.CategoryID,
		SubcategoryID: ps.SubcategoryID,
		SalePrice:     float64(ps.Price),
		IsSet:         true,
	}
	if err := s.itemRepo.Update(ctx, &item); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, ps); err != nil {
		return err
	}
	qty, err := s.calculateQuantity(ctx, ps)
	if err == nil {
		ps.Quantity = float64(qty)
		if err = s.itemRepo.SetStock(ctx, ps.ID, float64(qty)); err != nil {
			return err
		}
	}
	return nil
}

func (s *PriceSetService) DeletePriceSet(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return s.itemRepo.Delete(ctx, id)
}

func (s *PriceSetService) calculateQuantity(ctx context.Context, ps *models.PriceSet) (int, error) {
	qty := math.MaxInt32
	for _, it := range ps.Items {
		item, err := s.itemRepo.GetByID(ctx, it.ItemID)
		if err != nil {
			return 0, err
		}
		hours, err := s.isHoursCategory(ctx, item.CategoryID)
		if err != nil {
			return 0, err
		}
		if hours {
			continue
		}
		avail := int(item.Quantity / it.Quantity)
		if avail < qty {
			qty = avail
		}
	}
	if qty == math.MaxInt32 {
		qty = 0
	}
	return qty, nil
}
