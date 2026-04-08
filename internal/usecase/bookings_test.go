package usecase

import (
	"context"
	"testing"
	"time"

	"avito/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBooking_SlotInPast_ReturnsInvalidRequest(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()
	past := time.Now().Add(-time.Hour)

	uc := NewBookingService(
		&mockBookingRepo{},
		&mockSlotRepoBooking{slot: &entity.Slot{ID: slotID, Start: past}},
		NewConferenceUseCase(),
	)

	_, err := uc.Create(context.Background(), userID, slotID, false)
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeInvalidRequest, code)
}

func TestCreateBooking_SlotNotFound_Propagates(t *testing.T) {
	uc := NewBookingService(
		&mockBookingRepo{},
		&mockSlotRepoBooking{err: entity.New(entity.CodeSlotNotFound, "slot not found")},
		NewConferenceUseCase(),
	)

	_, err := uc.Create(context.Background(), uuid.New(), uuid.New(), false)
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeSlotNotFound, code)
}

func TestCreateBooking_Success_NoConference(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()
	future := time.Now().Add(time.Hour)

	uc := NewBookingService(
		&mockBookingRepo{},
		&mockSlotRepoBooking{slot: &entity.Slot{ID: slotID, Start: future}},
		NewConferenceUseCase(),
	)

	b, err := uc.Create(context.Background(), userID, slotID, false)
	require.NoError(t, err)
	assert.Equal(t, userID, b.UserID)
	assert.Equal(t, slotID, b.SlotID)
	assert.Equal(t, entity.BookingStatusActive, b.Status)
	assert.Nil(t, b.ConferenceLink)
}

func TestCreateBooking_WithConferenceLink(t *testing.T) {
	userID := uuid.New()
	slotID := uuid.New()
	future := time.Now().Add(time.Hour)

	uc := NewBookingService(
		&mockBookingRepo{},
		&mockSlotRepoBooking{slot: &entity.Slot{ID: slotID, Start: future}},
		NewConferenceUseCase(),
	)

	b, err := uc.Create(context.Background(), userID, slotID, true)
	require.NoError(t, err)
	require.NotNil(t, b.ConferenceLink)
	assert.Contains(t, *b.ConferenceLink, "https://example.com/meeting/")
}

func TestCancelBooking_AlreadyCancelled_Idempotent(t *testing.T) {
	bookingID := uuid.New()
	userID := uuid.New()
	existing := &entity.Booking{ID: bookingID, UserID: userID, Status: entity.BookingStatusCancelled}

	uc := NewBookingService(
		&mockBookingRepo{booking: existing},
		&mockSlotRepoBooking{},
		NewConferenceUseCase(),
	)

	b, err := uc.Cancel(context.Background(), bookingID, userID)
	require.NoError(t, err)
	assert.Equal(t, entity.BookingStatusCancelled, b.Status)
}

func TestCancelBooking_WrongUser_Forbidden(t *testing.T) {
	bookingID := uuid.New()
	owner := uuid.New()
	other := uuid.New()
	existing := &entity.Booking{ID: bookingID, UserID: owner, Status: entity.BookingStatusActive}

	uc := NewBookingService(
		&mockBookingRepo{booking: existing},
		&mockSlotRepoBooking{},
		NewConferenceUseCase(),
	)

	_, err := uc.Cancel(context.Background(), bookingID, other)
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeForbidden, code)
}

func TestCancelBooking_Success(t *testing.T) {
	bookingID := uuid.New()
	userID := uuid.New()
	existing := &entity.Booking{ID: bookingID, UserID: userID, Status: entity.BookingStatusActive}

	uc := NewBookingService(
		&mockBookingRepo{booking: existing},
		&mockSlotRepoBooking{},
		NewConferenceUseCase(),
	)

	b, err := uc.Cancel(context.Background(), bookingID, userID)
	require.NoError(t, err)
	assert.Equal(t, entity.BookingStatusCancelled, b.Status)
}

func TestCancelBooking_NotFound_PropagatesError(t *testing.T) {
	uc := NewBookingService(
		&mockBookingRepo{booking: nil, err: entity.New(entity.CodeBookingNotFound, "booking not found")},
		&mockSlotRepoBooking{},
		NewConferenceUseCase(),
	)

	_, err := uc.Cancel(context.Background(), uuid.New(), uuid.New())
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeBookingNotFound, code)
}

func TestBooking_GetByID(t *testing.T) {
	id := uuid.New()
	b := &entity.Booking{ID: id, Status: entity.BookingStatusActive}
	uc := NewBookingService(&mockBookingRepo{booking: b}, &mockSlotRepoBooking{}, NewConferenceUseCase())

	got, err := uc.GetByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
}

func TestBooking_ListAll(t *testing.T) {
	uc := NewBookingService(
		&mockBookingRepo{all: []entity.Booking{{ID: uuid.New()}}, total: 1},
		&mockSlotRepoBooking{},
		NewConferenceUseCase(),
	)
	list, total, err := uc.ListAll(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, 1, total)
}

func TestBooking_MyBookings(t *testing.T) {
	uid := uuid.New()
	uc := NewBookingService(
		&mockBookingRepo{all: []entity.Booking{{UserID: uid}}},
		&mockSlotRepoBooking{},
		NewConferenceUseCase(),
	)
	list, err := uc.MyBookings(context.Background(), uid)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}
