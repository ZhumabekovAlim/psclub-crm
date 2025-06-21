package services

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type BookingService struct {
	repo            *repositories.BookingRepository
	bookingItemRepo *repositories.BookingItemRepository
	clientRepo      *repositories.ClientRepository
	settingsRepo    *repositories.SettingsRepository
	priceItemRepo   *repositories.PriceItemRepository
	priceSetRepo    *repositories.PriceSetRepository
}

func NewBookingService(r *repositories.BookingRepository, itemRepo *repositories.BookingItemRepository, clientRepo *repositories.ClientRepository, settingsRepo *repositories.SettingsRepository, priceRepo *repositories.PriceItemRepository, setRepo *repositories.PriceSetRepository) *BookingService {
	return &BookingService{
		repo:            r,
		bookingItemRepo: itemRepo,
		clientRepo:      clientRepo,
		settingsRepo:    settingsRepo,
		priceItemRepo:   priceRepo,
		priceSetRepo:    setRepo,
	}
}

type stockChange struct {
	id  int
	qty int
}

func (s *BookingService) CreateBooking(ctx context.Context, b *models.Booking) (int, error) {
	// получить настройки для бонуса
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		log.Printf("settings get error: %v", err)
		return 0, err
	}
	if err := s.checkStock(ctx, b.Items); err != nil {
		log.Printf("check stock error: %v", err)
		return 0, err
	}
	id, err := s.repo.CreateWithItems(ctx, b)
	if err != nil {
		log.Printf("repository create error: %v", err)
		return 0, err
	}
	if err := s.decreaseStock(ctx, b.Items); err != nil {
		_ = s.repo.Delete(ctx, id)
		log.Printf("decrease stock error: %v", err)
		return 0, err
	}
	// Списываем использованные бонусы
	if b.BonusUsed > 0 {
		_ = s.clientRepo.AddBonus(ctx, b.ClientID, -b.BonusUsed)
	}
	// Начисляем бонусы с суммы, оплаченной деньгами
	paid := b.TotalAmount - b.BonusUsed
	if paid < 0 {
		paid = 0
	}
	bonus := int(float64(paid) * float64(settings.BonusPercent) / 100)
	_ = s.clientRepo.AddBonus(ctx, b.ClientID, bonus)
	_ = s.clientRepo.AddVisits(ctx, b.ClientID, 1)
	_ = s.clientRepo.AddIncome(ctx, b.ClientID, b.TotalAmount)
	return id, nil
}

func (s *BookingService) GetAllBookings(ctx context.Context) ([]models.Booking, error) {
	bookings, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	// загрузить позиции для каждой брони
	for i := range bookings {
		items, _ := s.bookingItemRepo.GetByBookingID(ctx, bookings[i].ID)
		bookings[i].Items = items
	}
	return bookings, nil
}

func (s *BookingService) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	items, _ := s.bookingItemRepo.GetByBookingID(ctx, b.ID)
	b.Items = items
	return b, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, b *models.Booking) error {
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return err
	}
	current, err := s.repo.GetByID(ctx, b.ID)
	if err != nil {
		return err
	}
	limit := current.EndTime.Add(time.Duration(settings.BlockTime) * time.Minute)
	if time.Now().After(limit) {
		return errors.New("booking can no longer be modified")
	}

	if err := s.checkStock(ctx, b.Items); err != nil {
		return err
	}
	if err := s.decreaseStock(ctx, b.Items); err != nil {
		return err
	}
	if err := s.repo.Update(ctx, b); err != nil {
		// rollback stock on failure
		s.increaseStock(ctx, b.Items)
		return err
	}
	return nil

}

func (s *BookingService) DeleteBooking(ctx context.Context, id int) error {
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return err
	}
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	limit := b.EndTime.Add(time.Duration(settings.BlockTime) * time.Minute)
	if time.Now().After(limit) {
		return errors.New("booking can no longer be removed")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// отменяем начисленные бонусы и возвращаем использованные
	paid := b.TotalAmount - b.BonusUsed
	if paid < 0 {
		paid = 0
	}
	bonus := int(float64(paid) * float64(settings.BonusPercent) / 100)
	_ = s.clientRepo.AddBonus(ctx, b.ClientID, -bonus)
	if b.BonusUsed > 0 {
		_ = s.clientRepo.AddBonus(ctx, b.ClientID, b.BonusUsed)
	}
	_ = s.clientRepo.AddVisits(ctx, b.ClientID, -1)
	_ = s.clientRepo.AddIncome(ctx, b.ClientID, -b.TotalAmount)
	return nil
}

func (s *BookingService) decreaseStock(ctx context.Context, items []models.BookingItem) error {

	var changes []stockChange

	affected := make(map[int]struct{})
	for _, it := range items {
		if it.Quantity <= 0 {
			continue
		}

		pi, err := s.priceItemRepo.GetByID(ctx, it.ItemID)
		if err != nil {
			s.restoreChanges(ctx, changes)
			return err
		}
		if err := s.priceItemRepo.DecreaseStock(ctx, it.ItemID, it.Quantity); err != nil {
			s.restoreChanges(ctx, changes)
			return err
		}

		changes = append(changes, stockChange{id: it.ItemID, qty: it.Quantity})

		affected[it.ItemID] = struct{}{}
		if pi.IsSet {
			set, err := s.priceSetRepo.GetByID(ctx, pi.ID)
			if err != nil {

				s.restoreChanges(ctx, changes)

				return err
			}
			for _, si := range set.Items {
				if err := s.priceItemRepo.DecreaseStock(ctx, si.ItemID, si.Quantity*it.Quantity); err != nil {
					s.restoreChanges(ctx, changes)
					return err
				}

				changes = append(changes, stockChange{id: si.ItemID, qty: si.Quantity * it.Quantity})

				affected[si.ItemID] = struct{}{}
			}
		}
	}

	if err := s.updateSetQuantities(ctx, affected); err != nil {
		s.restoreChanges(ctx, changes)
		return err
	}
	return nil

}

func (s *BookingService) updateSetQuantities(ctx context.Context, affected map[int]struct{}) error {
	updated := make(map[int]struct{})
	for itemID := range affected {
		sets, err := s.priceSetRepo.GetByItem(ctx, itemID)
		if err != nil {
			return err
		}
		for _, set := range sets {
			if _, ok := updated[set.ID]; ok {
				continue
			}
			qty := math.MaxInt32
			for _, si := range set.Items {
				it, err := s.priceItemRepo.GetByID(ctx, si.ItemID)
				if err != nil {
					return err
				}
				avail := it.Quantity / si.Quantity
				if avail < qty {
					qty = avail
				}
			}
			if qty == math.MaxInt32 {
				qty = 0
			}
			if err := s.priceItemRepo.SetStock(ctx, set.ID, qty); err != nil {
				return err
			}
			updated[set.ID] = struct{}{}
		}
	}
	return nil
}

func (s *BookingService) checkStock(ctx context.Context, items []models.BookingItem) error {
	for _, it := range items {
		pi, err := s.priceItemRepo.GetByID(ctx, it.ItemID)
		if err != nil {
			return err
		}
		if pi.Quantity < it.Quantity {
			return errors.New("insufficient stock")
		}
		if pi.IsSet {
			set, err := s.priceSetRepo.GetByID(ctx, pi.ID)
			if err != nil {
				return err
			}
			for _, si := range set.Items {
				sub, err := s.priceItemRepo.GetByID(ctx, si.ItemID)
				if err != nil {
					return err
				}
				if sub.Quantity < si.Quantity*it.Quantity {
					return errors.New("insufficient stock")
				}
			}
		}
	}
	return nil
}

func (s *BookingService) restoreChanges(ctx context.Context, changes []stockChange) {

	affected := make(map[int]struct{})
	for _, c := range changes {
		_ = s.priceItemRepo.IncreaseStock(ctx, c.id, c.qty)
		affected[c.id] = struct{}{}
	}
	_ = s.updateSetQuantities(ctx, affected)
}

func (s *BookingService) increaseStock(ctx context.Context, items []models.BookingItem) {

	var changes []stockChange

	for _, it := range items {
		pi, err := s.priceItemRepo.GetByID(ctx, it.ItemID)
		if err != nil {
			continue
		}
		_ = s.priceItemRepo.IncreaseStock(ctx, it.ItemID, it.Quantity)
		changes = append(changes, stockChange{id: it.ItemID, qty: it.Quantity})

		if pi.IsSet {
			set, err := s.priceSetRepo.GetByID(ctx, pi.ID)
			if err != nil {
				continue
			}
			for _, si := range set.Items {
				_ = s.priceItemRepo.IncreaseStock(ctx, si.ItemID, si.Quantity*it.Quantity)

				changes = append(changes, stockChange{id: si.ItemID, qty: si.Quantity * it.Quantity})

			}
		}
	}
	s.restoreChanges(ctx, changes) // update sets after restocking
}
