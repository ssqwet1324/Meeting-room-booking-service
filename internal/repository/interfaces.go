package repository

import (
	"avito/internal/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

// UserRepository описывает интерфейс репозитория пользователей.
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, email, passwordHash string, role entity.Role) (*entity.User, error)
}

// RoomRepository описывает интерфейс репозитория комнат.
type RoomRepository interface {
	ListRooms(ctx context.Context) ([]entity.Room, error)
	GetRoomByID(ctx context.Context, id uuid.UUID) (*entity.Room, error)
	CreateRoom(ctx context.Context, newRoom entity.NewRoom) (*entity.Room, error)
}

// ScheduleRepository описывает интерфейс репозитория расписаний.
type ScheduleRepository interface {
	GetScheduleByRoomID(ctx context.Context, roomID uuid.UUID) (*entity.Schedule, error)
	CreateSchedule(ctx context.Context, schedule entity.NewSchedule) (*entity.Schedule, error)
}

// SlotRepository описывает интерфейс репозитория слотов.
type SlotRepository interface {
	GetSlotByID(ctx context.Context, id uuid.UUID) (*entity.Slot, error)
	GetByRoomAndDate(ctx context.Context, roomID uuid.UUID, dayStart, dayEnd time.Time) ([]entity.Slot, error)
	GetAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, dayStart, dayEnd time.Time) ([]entity.Slot, error)
	InsertNewSlotsFromRoom(ctx context.Context, slots []entity.Slot) error
}

// BookingRepository описывает интерфейс репозитория бронирований.
type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *entity.Booking) (*entity.Booking, error)
	GetBookingByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error)
	GetBookingByUserID(ctx context.Context, userID uuid.UUID, from time.Time) ([]entity.Booking, error)
	ListAll(ctx context.Context, offset, limit int) ([]entity.Booking, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.BookingStatus) (*entity.Booking, error)
}
