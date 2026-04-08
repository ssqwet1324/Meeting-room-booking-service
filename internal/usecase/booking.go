package usecase

import (
	"avito/internal/entity"
	"avito/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

// BookingUseCase представляет use case бронирований.
type BookingUseCase struct {
	bookings   repository.BookingRepository
	slots      repository.SlotRepository
	conference *ConferenceUseCase
}

// NewBookingService создает use case бронирований.
func NewBookingService(bookings repository.BookingRepository, slots repository.SlotRepository, conference *ConferenceUseCase) *BookingUseCase {
	return &BookingUseCase{
		bookings:   bookings,
		slots:      slots,
		conference: conference,
	}
}

// Create - забронировать слот
func (u *BookingUseCase) Create(ctx context.Context, userID uuid.UUID, slotID uuid.UUID, createConferenceLink bool) (*entity.Booking, error) {
	// проверка существования слота
	slotExists, err := u.slots.GetSlotByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	// проверяем не в прошлом ли время, на которое бронируем
	if !slotExists.Start.After(time.Now()) {
		return nil, entity.New(entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
	}

	booking := &entity.Booking{
		SlotID: slotID,
		UserID: userID,
		Status: entity.BookingStatusActive,
	}

	// берем ссылку на конференцию
	if createConferenceLink {
		id := uuid.New()
		link, err := u.conference.CreateURL(id)
		if err == nil {
			booking.ConferenceLink = &link
		}
	}

	crated, err := u.bookings.CreateBooking(ctx, booking)
	if err != nil {
		return nil, err
	}

	return crated, nil
}

// Cancel - отменить бронь
func (u *BookingUseCase) Cancel(ctx context.Context, bookingID, userID uuid.UUID) (*entity.Booking, error) {
	booking, err := u.bookings.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.Status == entity.BookingStatusCancelled {
		return booking, nil
	}

	if booking.UserID != userID {
		return nil, entity.New(entity.CodeForbidden, entity.ErrMsgCancelOtherUser)
	}

	return u.bookings.UpdateStatus(ctx, bookingID, entity.BookingStatusCancelled)
}

// GetByID - получить бронь по id
func (u *BookingUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	return u.bookings.GetBookingByID(ctx, id)
}

// ListAll - получить список броней с пагинацией
func (u *BookingUseCase) ListAll(ctx context.Context, page, pageSize int) ([]entity.Booking, int, error) {
	offset := (page - 1) * pageSize

	return u.bookings.ListAll(ctx, offset, pageSize)
}

// MyBookings - получить мои брони
func (u *BookingUseCase) MyBookings(ctx context.Context, userID uuid.UUID) ([]entity.Booking, error) {
	now := time.Now()

	return u.bookings.GetBookingByUserID(ctx, userID, now)
}
