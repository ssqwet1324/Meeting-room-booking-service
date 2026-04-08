package usecase

import (
	"context"
	"time"

	"avito/internal/entity"

	"github.com/google/uuid"
)

// --- user (auth) ---

type mockUserRepo struct {
	user *entity.User
	err  error
}

func (m *mockUserRepo) GetByID(_ context.Context, _ uuid.UUID) (*entity.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) GetByEmail(_ context.Context, _ string) (*entity.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) Create(_ context.Context, email, hash string, role entity.Role) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	now := time.Now()
	return &entity.User{ID: uuid.New(), Email: email, PasswordHash: hash, Role: role, CreatedAt: &now}, nil
}

// --- booking flow ---

type mockBookingRepo struct {
	booking *entity.Booking
	all     []entity.Booking
	total   int
	err     error
}

func (m *mockBookingRepo) CreateBooking(_ context.Context, b *entity.Booking) (*entity.Booking, error) {
	if m.err != nil {
		return nil, m.err
	}
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return b, nil
}

func (m *mockBookingRepo) GetBookingByID(_ context.Context, _ uuid.UUID) (*entity.Booking, error) {
	return m.booking, m.err
}

func (m *mockBookingRepo) GetBookingByUserID(_ context.Context, _ uuid.UUID, _ time.Time) ([]entity.Booking, error) {
	return m.all, m.err
}

func (m *mockBookingRepo) ListAll(_ context.Context, _, _ int) ([]entity.Booking, int, error) {
	return m.all, m.total, m.err
}

func (m *mockBookingRepo) UpdateStatus(_ context.Context, _ uuid.UUID, status entity.BookingStatus) (*entity.Booking, error) {
	if m.booking != nil {
		m.booking.Status = status
		return m.booking, m.err
	}
	return nil, m.err
}

type mockSlotRepoBooking struct {
	slot *entity.Slot
	err  error
}

func (m *mockSlotRepoBooking) GetSlotByID(_ context.Context, _ uuid.UUID) (*entity.Slot, error) {
	return m.slot, m.err
}

func (m *mockSlotRepoBooking) GetByRoomAndDate(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]entity.Slot, error) {
	return nil, nil
}

func (m *mockSlotRepoBooking) GetAvailableByRoomAndDate(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]entity.Slot, error) {
	return nil, nil
}

func (m *mockSlotRepoBooking) InsertNewSlotsFromRoom(_ context.Context, _ []entity.Slot) error {
	return nil
}

// --- slots / schedule / rooms ---

type mockRoomRepo struct {
	room  *entity.Room
	rooms []entity.Room
	err   error
}

func (m *mockRoomRepo) ListRooms(_ context.Context) ([]entity.Room, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rooms, nil
}

func (m *mockRoomRepo) GetRoomByID(_ context.Context, _ uuid.UUID) (*entity.Room, error) {
	return m.room, m.err
}

func (m *mockRoomRepo) CreateRoom(_ context.Context, newRoom entity.NewRoom) (*entity.Room, error) {
	if m.err != nil {
		return nil, m.err
	}
	now := time.Now()
	return &entity.Room{ID: uuid.New(), Name: newRoom.Name, CreatedAt: &now}, nil
}

type mockScheduleRepo struct {
	schedule *entity.Schedule
	err      error
}

func (m *mockScheduleRepo) GetScheduleByRoomID(_ context.Context, _ uuid.UUID) (*entity.Schedule, error) {
	return m.schedule, m.err
}

func (m *mockScheduleRepo) CreateSchedule(_ context.Context, s entity.NewSchedule) (*entity.Schedule, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &entity.Schedule{
		ID:         uuid.New(),
		RoomID:     s.RoomID,
		DaysOfWeek: s.DaysOfWeek,
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
	}, nil
}

type mockSlotRepoSlots struct {
	existing  []entity.Slot
	available []entity.Slot
	inserted  []entity.Slot
	slotByID  *entity.Slot
	err       error
}

func (m *mockSlotRepoSlots) GetSlotByID(_ context.Context, _ uuid.UUID) (*entity.Slot, error) {
	if m.slotByID != nil {
		return m.slotByID, m.err
	}
	return nil, m.err
}

func (m *mockSlotRepoSlots) GetByRoomAndDate(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]entity.Slot, error) {
	return m.existing, m.err
}

func (m *mockSlotRepoSlots) GetAvailableByRoomAndDate(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]entity.Slot, error) {
	return m.available, m.err
}

func (m *mockSlotRepoSlots) InsertNewSlotsFromRoom(_ context.Context, slots []entity.Slot) error {
	m.inserted = slots
	return nil
}
