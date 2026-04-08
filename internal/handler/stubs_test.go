package handler_test

import (
	"context"
	"testing"
	"time"

	"avito/internal/entity"
	"avito/internal/handler"
	"avito/internal/usecase"

	"github.com/google/uuid"
)

type stubUserRepo struct {
	user *entity.User
	err  error
}

func (s *stubUserRepo) GetByID(_ context.Context, _ uuid.UUID) (*entity.User, error) {
	return s.user, s.err
}

func (s *stubUserRepo) GetByEmail(_ context.Context, _ string) (*entity.User, error) {
	return s.user, s.err
}

func (s *stubUserRepo) Create(_ context.Context, _, _ string, _ entity.Role) (*entity.User, error) {
	if s.err != nil {
		return nil, s.err
	}
	now := time.Now()
	return &entity.User{ID: uuid.New(), Email: "x@y.com", Role: entity.RoleUser, CreatedAt: &now}, nil
}

type stubBookingRepo struct {
	booking *entity.Booking
	all     []entity.Booking
	total   int
	err     error
}

func (s *stubBookingRepo) CreateBooking(_ context.Context, b *entity.Booking) (*entity.Booking, error) {
	if s.err != nil {
		return nil, s.err
	}
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return b, nil
}

func (s *stubBookingRepo) GetBookingByID(_ context.Context, _ uuid.UUID) (*entity.Booking, error) {
	return s.booking, s.err
}

func (s *stubBookingRepo) GetBookingByUserID(_ context.Context, _ uuid.UUID, _ time.Time) ([]entity.Booking, error) {
	return s.all, s.err
}

func (s *stubBookingRepo) ListAll(_ context.Context, _, _ int) ([]entity.Booking, int, error) {
	return s.all, s.total, s.err
}

func (s *stubBookingRepo) UpdateStatus(_ context.Context, _ uuid.UUID, status entity.BookingStatus) (*entity.Booking, error) {
	if s.booking != nil {
		s.booking.Status = status
		return s.booking, s.err
	}
	return nil, s.err
}

type stubSlotRepoBooking struct {
	slot *entity.Slot
	err  error
}

func (s *stubSlotRepoBooking) GetSlotByID(_ context.Context, _ uuid.UUID) (*entity.Slot, error) {
	return s.slot, s.err
}

func (s *stubSlotRepoBooking) GetByRoomAndDate(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]entity.Slot, error) {
	return nil, nil
}

func (s *stubSlotRepoBooking) GetAvailableByRoomAndDate(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]entity.Slot, error) {
	return nil, nil
}

func (s *stubSlotRepoBooking) InsertNewSlotsFromRoom(_ context.Context, _ []entity.Slot) error {
	return nil
}

type stubRoomRepo struct {
	rooms     []entity.Room
	room      *entity.Room
	roomsErr  error
	roomErr   error
	createErr error
}

func (s *stubRoomRepo) ListRooms(_ context.Context) ([]entity.Room, error) {
	if s.roomsErr != nil {
		return nil, s.roomsErr
	}
	return s.rooms, nil
}

func (s *stubRoomRepo) GetRoomByID(_ context.Context, id uuid.UUID) (*entity.Room, error) {
	if s.roomErr != nil {
		return nil, s.roomErr
	}
	if s.room != nil {
		return s.room, nil
	}
	return &entity.Room{ID: id, Name: "room"}, nil
}

func (s *stubRoomRepo) CreateRoom(_ context.Context, newRoom entity.NewRoom) (*entity.Room, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	now := time.Now()
	return &entity.Room{ID: uuid.New(), Name: newRoom.Name, CreatedAt: &now}, nil
}

type stubScheduleRepo struct {
	sch       *entity.Schedule
	err       error
	createErr error
}

func (s *stubScheduleRepo) GetScheduleByRoomID(_ context.Context, _ uuid.UUID) (*entity.Schedule, error) {
	return s.sch, s.err
}

func (s *stubScheduleRepo) CreateSchedule(_ context.Context, ns entity.NewSchedule) (*entity.Schedule, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	return &entity.Schedule{
		ID:         uuid.New(),
		RoomID:     ns.RoomID,
		DaysOfWeek: ns.DaysOfWeek,
		StartTime:  ns.StartTime,
		EndTime:    ns.EndTime,
	}, nil
}

type stubSlotRepoFull struct {
	existing  []entity.Slot
	available []entity.Slot
	inserted  []entity.Slot
	err       error
}

func (s *stubSlotRepoFull) GetSlotByID(_ context.Context, _ uuid.UUID) (*entity.Slot, error) {
	return nil, nil
}

func (s *stubSlotRepoFull) GetByRoomAndDate(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]entity.Slot, error) {
	return s.existing, s.err
}

func (s *stubSlotRepoFull) GetAvailableByRoomAndDate(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]entity.Slot, error) {
	return s.available, s.err
}

func (s *stubSlotRepoFull) InsertNewSlotsFromRoom(_ context.Context, slots []entity.Slot) error {
	s.inserted = slots
	return nil
}

func newTestHandler(t *testing.T) (*handler.Handler, *usecase.AuthUseCase) {
	t.Helper()
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "secret")
	bookings := usecase.NewBookingService(&stubBookingRepo{}, &stubSlotRepoBooking{}, usecase.NewConferenceUseCase())
	h := handler.New(auth, &usecase.RoomUseCase{}, &usecase.ScheduleUseCase{}, &usecase.SlotUseCase{}, bookings)
	return h, auth
}

func newHandlerWithRooms(t *testing.T, roomID uuid.UUID) *handler.Handler {
	t.Helper()
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "secret")
	roomUC := usecase.NewRoomUseCase(&stubRoomRepo{
		rooms: []entity.Room{{ID: roomID, Name: "Hall"}},
		room:  &entity.Room{ID: roomID, Name: "Hall"},
	})
	scheduleUC := usecase.NewScheduleUseCase(
		&stubScheduleRepo{err: entity.New(entity.CodeNotFound, "no schedule")},
		&stubRoomRepo{room: &entity.Room{ID: roomID}},
	)
	slotUC := usecase.NewSlotUseCase(
		&stubSlotRepoFull{},
		&stubScheduleRepo{err: entity.New(entity.CodeNotFound, "no schedule")},
		&stubRoomRepo{room: &entity.Room{ID: roomID}},
	)
	bookings := usecase.NewBookingService(&stubBookingRepo{}, &stubSlotRepoBooking{}, usecase.NewConferenceUseCase())
	return handler.New(auth, roomUC, scheduleUC, slotUC, bookings)
}
