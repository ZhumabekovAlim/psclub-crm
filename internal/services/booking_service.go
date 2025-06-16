package services

import (
	"context"
	"errors"
	"time"

	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type BookingService struct {
	repo            *repositories.BookingRepository
	bookingItemRepo *repositories.BookingItemRepository
	clientRepo      *repositories.ClientRepository
	settingsRepo    *repositories.SettingsRepository
}

func NewBookingService(r *repositories.BookingRepository, itemRepo *repositories.BookingItemRepository, clientRepo *repositories.ClientRepository, settingsRepo *repositories.SettingsRepository) *BookingService {
	return &BookingService{
		repo:            r,
		bookingItemRepo: itemRepo,
		clientRepo:      clientRepo,
		settingsRepo:    settingsRepo,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, b *models.Booking) (int, error) {
	// получить настройки для бонуса
	settings, err := s.settingsRepo.Get(ctx)
	if err != nil {
		return 0, err
	}
	id, err := s.repo.CreateWithItems(ctx, b)
	if err != nil {
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
	limit := current.StartTime.Add(time.Duration(settings.BlockTime) * time.Minute)
	if time.Now().After(limit) {
		return errors.New("booking can no longer be modified")
	}
	return s.repo.Update(ctx, b)
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
	limit := b.StartTime.Add(time.Duration(settings.BlockTime) * time.Minute)
	if time.Now().After(limit) {
		return errors.New("booking can no longer be removed")
	}
	return s.repo.Delete(ctx, id)
}
