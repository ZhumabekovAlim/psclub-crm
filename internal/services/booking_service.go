package services

import (
	"context"
	"errors"
	"log"
	"math"
	"strings"
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
	categoryRepo    *repositories.CategoryRepository
	paymentRepo     *repositories.BookingPaymentRepository
	paymentTypeRepo *repositories.PaymentTypeRepository
	cashboxService  *CashboxService
}

func NewBookingService(r *repositories.BookingRepository, itemRepo *repositories.BookingItemRepository, clientRepo *repositories.ClientRepository, settingsRepo *repositories.SettingsRepository, priceRepo *repositories.PriceItemRepository, setRepo *repositories.PriceSetRepository, categoryRepo *repositories.CategoryRepository, paymentRepo *repositories.BookingPaymentRepository, ptRepo *repositories.PaymentTypeRepository, cbService *CashboxService) *BookingService {
	return &BookingService{
		repo:            r,
		bookingItemRepo: itemRepo,
		clientRepo:      clientRepo,
		settingsRepo:    settingsRepo,
		priceItemRepo:   priceRepo,
		priceSetRepo:    setRepo,
		categoryRepo:    categoryRepo,
		paymentRepo:     paymentRepo,
		paymentTypeRepo: ptRepo,
		cashboxService:  cbService,
	}
}

type stockChange struct {
	id  int
	qty float64
}

func bookingsEqual(old *models.Booking, newB *models.Booking, oldItems []models.BookingItem) bool {
	if old.ClientID != newB.ClientID ||
		old.TableID != newB.TableID ||
		old.UserID != newB.UserID ||
		!old.StartTime.Equal(newB.StartTime) ||
		!old.EndTime.Equal(newB.EndTime) ||
		old.Note != newB.Note ||
		old.Discount != newB.Discount ||
		old.DiscountReason != newB.DiscountReason ||
		old.TotalAmount != newB.TotalAmount ||
		old.BonusUsed != newB.BonusUsed ||
		old.PaymentStatus != newB.PaymentStatus ||
		old.PaymentTypeID != newB.PaymentTypeID {
		return false
	}

	if len(oldItems) != len(newB.Items) {
		return false
	}

	m := make(map[int]models.BookingItem)
	for _, it := range oldItems {
		m[it.ItemID] = it
	}
	for _, it := range newB.Items {
		o, ok := m[it.ItemID]
		if !ok || o.Quantity != it.Quantity || o.Price != it.Price || o.Discount != it.Discount {
			return false
		}
	}
	return true
}

func (s *BookingService) isHoursCategory(ctx context.Context, categoryID int) (bool, error) {
	cat, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return false, err
	}
	return cat.Name == hoursCategoryName, nil
}

func (s *BookingService) CreateBooking(ctx context.Context, b *models.Booking) (int, error) {
	if len(b.Payments) > 0 {
		b.PaymentTypeID = b.Payments[0].PaymentTypeID
	}
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
	if b.ClientID > 0 {
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
	}

	if err := s.paymentRepo.Create(ctx, id, b.Payments); err != nil {
		_ = s.repo.Delete(ctx, id)
		return 0, err
	}

	if strings.ToLower(b.PaymentStatus) == "paid" && s.cashboxService != nil {
		for _, p := range b.Payments {
			pt, err := s.paymentTypeRepo.GetByID(ctx, p.PaymentTypeID)
			if err != nil {
				continue
			}
			name := strings.ToLower(pt.Name)
			if strings.Contains(name, "наличными") {
				_ = s.cashboxService.AddIncome(ctx, float64(p.Amount))
			}
		}
	}
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
		for j := range items {
			if items[j].Quantity != 0 {
				items[j].ItemPrice = float64(items[j].Price) / items[j].Quantity
			}
		}
		bookings[i].Items = items
		pays, _ := s.paymentRepo.GetByBookingID(ctx, bookings[i].ID)
		bookings[i].Payments = pays
		if len(pays) > 0 {
			bookings[i].PaymentType = pays[0].PaymentType
		}
	}
	return bookings, nil
}

func (s *BookingService) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	items, _ := s.bookingItemRepo.GetByBookingID(ctx, b.ID)
	for i := range items {
		if items[i].Quantity != 0 {
			items[i].ItemPrice = float64(items[i].Price) / items[i].Quantity
		}
	}
	b.Items = items
	pays, _ := s.paymentRepo.GetByBookingID(ctx, b.ID)
	b.Payments = pays
	if len(pays) > 0 {
		b.PaymentType = pays[0].PaymentType
	}
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
	currentItems, _ := s.bookingItemRepo.GetByBookingID(ctx, b.ID)
	currentPays, _ := s.paymentRepo.GetByBookingID(ctx, b.ID)

	equal := bookingsEqual(current, b, currentItems)
	if equal && len(currentPays) == len(b.Payments) {
		eq := true
		for i, p := range b.Payments {
			if p.PaymentTypeID != currentPays[i].PaymentTypeID || p.Amount != currentPays[i].Amount {
				eq = false
				break
			}
		}
		if eq {
			return nil
		}
	}

	if len(b.Payments) > 0 {
		b.PaymentTypeID = b.Payments[0].PaymentTypeID
	}
	limit := current.EndTime.Add(time.Duration(settings.BlockTime) * time.Minute)
	if time.Now().After(limit) {
		return errors.New("booking can no longer be modified")
	}

	s.increaseStock(ctx, currentItems)

	if err := s.checkStock(ctx, b.Items); err != nil {
		err := s.decreaseStock(ctx, currentItems)
		if err != nil {
			return err
		}
		return err
	}
	if err := s.decreaseStock(ctx, b.Items); err != nil {
		err := s.decreaseStock(ctx, currentItems)
		if err != nil {
			return err
		}
		return err
	}
	oldCash := 0.0
	newCash := 0.0
	if s.cashboxService != nil {
		oldCash = s.getCashAmount(ctx, current)
		newCash = s.getCashAmount(ctx, b)
	}

	if err := s.repo.UpdateWithItems(ctx, b); err != nil {
		// rollback stock on failure
		s.increaseStock(ctx, b.Items)
		err := s.decreaseStock(ctx, currentItems)
		if err != nil {
			return err
		}
		return err
	}
	_ = s.paymentRepo.DeleteByBookingID(ctx, b.ID)
	_ = s.paymentRepo.Create(ctx, b.ID, b.Payments)
	if s.cashboxService != nil {
		diff := newCash - oldCash
		if diff > 0 {
			_ = s.cashboxService.AddIncome(ctx, diff)
		} else if diff < 0 {
			_ = s.cashboxService.RemoveIncome(ctx, -diff)
		}
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
	cash := 0.0
	if s.cashboxService != nil {
		cash = s.getCashAmount(ctx, b)
	}
	items, err := s.bookingItemRepo.GetByBookingID(ctx, id)
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
	_ = s.paymentRepo.DeleteByBookingID(ctx, id)
	s.increaseStock(ctx, items)

	if b.ClientID > 0 {
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
	}
	if s.cashboxService != nil && cash > 0 {
		_ = s.cashboxService.RemoveIncome(ctx, cash)
	}
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
		isHours, err := s.isHoursCategory(ctx, pi.CategoryID)
		if err != nil {
			s.restoreChanges(ctx, changes)
			return err
		}
		if !isHours {
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
					sub, err := s.priceItemRepo.GetByID(ctx, si.ItemID)
					if err != nil {
						s.restoreChanges(ctx, changes)
						return err
					}
					hoursSub, err := s.isHoursCategory(ctx, sub.CategoryID)
					if err != nil {
						s.restoreChanges(ctx, changes)
						return err
					}
					if hoursSub {
						continue
					}
					if err := s.priceItemRepo.DecreaseStock(ctx, si.ItemID, si.Quantity*it.Quantity); err != nil {
						s.restoreChanges(ctx, changes)
						return err
					}

					changes = append(changes, stockChange{id: si.ItemID, qty: si.Quantity * it.Quantity})

					affected[si.ItemID] = struct{}{}
				}
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
				hours, err := s.isHoursCategory(ctx, it.CategoryID)
				if err != nil {
					return err
				}
				if hours {
					continue
				}
				avail := int(it.Quantity / si.Quantity)
				if avail < qty {
					qty = avail
				}
			}
			if qty == math.MaxInt32 {
				qty = 0
			}
			if err := s.priceItemRepo.SetStock(ctx, set.ID, float64(qty)); err != nil {
				return err
			}
			updated[set.ID] = struct{}{}
		}
	}
	return nil
}

// calculateSetQuantity determines the maximum number of sets that can be
// assembled based on the stock levels of the items included in the set.
func (s *BookingService) calculateSetQuantity(ctx context.Context, ps *models.PriceSet) (float64, error) {
	qty := math.MaxFloat64
	for _, it := range ps.Items {
		item, err := s.priceItemRepo.GetByID(ctx, it.ItemID)
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
		avail := item.Quantity / it.Quantity
		if avail < qty {
			qty = avail
		}
	}
	if qty == math.MaxFloat64 {
		qty = 0
	}
	return qty, nil
}

func (s *BookingService) checkStock(ctx context.Context, items []models.BookingItem) error {
	for _, it := range items {
		pi, err := s.priceItemRepo.GetByID(ctx, it.ItemID)
		if err != nil {
			return err
		}

		isHours, err := s.isHoursCategory(ctx, pi.CategoryID)
		if err != nil {
			return err
		}
		if isHours {
			continue
		}

		if pi.IsSet {
			set, err := s.priceSetRepo.GetByID(ctx, pi.ID)
			if err != nil {
				return err
			}
			qty, err := s.calculateSetQuantity(ctx, set)
			if err != nil {
				return err
			}
			if err := s.priceItemRepo.SetStock(ctx, set.ID, float64(qty)); err != nil {
				return err
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
	affected := make(map[int]struct{})

	for _, it := range items {
		pi, err := s.priceItemRepo.GetByID(ctx, it.ItemID)
		if err != nil {
			continue
		}
		hours, err := s.isHoursCategory(ctx, pi.CategoryID)
		if err != nil {
			continue
		}
		if hours {
			continue
		}

		_ = s.priceItemRepo.IncreaseStock(ctx, it.ItemID, it.Quantity)
		affected[it.ItemID] = struct{}{}

		if pi.IsSet {
			set, err := s.priceSetRepo.GetByID(ctx, pi.ID)
			if err != nil {
				continue
			}
			for _, si := range set.Items {
				sub, err := s.priceItemRepo.GetByID(ctx, si.ItemID)
				if err != nil {
					continue
				}
				hoursSub, err := s.isHoursCategory(ctx, sub.CategoryID)
				if err != nil || hoursSub {
					continue
				}
				_ = s.priceItemRepo.IncreaseStock(ctx, si.ItemID, si.Quantity*it.Quantity)
				affected[si.ItemID] = struct{}{}
			}
		}
	}

	_ = s.updateSetQuantities(ctx, affected)
}

func (s *BookingService) getCashAmount(ctx context.Context, b *models.Booking) float64 {
	if strings.ToLower(b.PaymentStatus) != "paid" {
		return 0
	}
	total := 0.0
	for _, p := range b.Payments {
		pt, err := s.paymentTypeRepo.GetByID(ctx, p.PaymentTypeID)
		if err != nil {
			continue
		}
		name := strings.ToLower(pt.Name)
		if strings.Contains(name, "наличными") {
			total += float64(p.Amount)
		}
	}
	return total
}
