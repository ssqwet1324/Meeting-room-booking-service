package usecase

import (
	"avito/internal/entity"
	"avito/internal/repository"
	"context"

	"github.com/google/uuid"
)

// RoomUseCase представляет use case комнат.
type RoomUseCase struct {
	rooms repository.RoomRepository
}

// NewRoomUseCase создает use case комнат.
func NewRoomUseCase(rooms repository.RoomRepository) *RoomUseCase {
	return &RoomUseCase{rooms: rooms}
}

// CreateRoom - создание комнаты
func (s *RoomUseCase) CreateRoom(ctx context.Context, newRoom entity.NewRoom) (*entity.Room, error) {
	return s.rooms.CreateRoom(ctx, newRoom)
}

// GetRoomByID - получение комнаты по id
func (s *RoomUseCase) GetRoomByID(ctx context.Context, id uuid.UUID) (*entity.Room, error) {
	return s.rooms.GetRoomByID(ctx, id)
}

// ListRooms - получение списка комнат
func (s *RoomUseCase) ListRooms(ctx context.Context) ([]entity.Room, error) {
	return s.rooms.ListRooms(ctx)
}
