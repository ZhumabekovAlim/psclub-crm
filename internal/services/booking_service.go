package services

import (
	"context"
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
	bonus := int(float64(b.TotalAmount) * float64(settings.BonusPercent) / 100)
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
	return s.repo.Update(ctx, b)
}

func (s *BookingService) DeleteBooking(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
