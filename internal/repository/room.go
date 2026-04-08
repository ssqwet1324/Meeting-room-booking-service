package repository

import (
	"avito/internal/entity"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type roomRepository struct {
	db *pgxpool.Pool
}

// NewRoomRepository создает репозиторий комнат.
func NewRoomRepository(db *pgxpool.Pool) RoomRepository {
	return &roomRepository{db: db}
}

// CreateRoom - создание комнаты
func (repo *roomRepository) CreateRoom(ctx context.Context, newRoom entity.NewRoom) (*entity.Room, error) {
	row := repo.db.QueryRow(ctx,
		`INSERT INTO rooms (name, description, capacity) VALUES ($1, $2, $3)
		 RETURNING id, name, description, capacity, created_at`,
		newRoom.Name, newRoom.Description, newRoom.Capacity)

	var room entity.Room
	if err := row.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt); err != nil {
		slog.Error("Repository: Error create room",
			slog.Any("room", newRoom),
			slog.Any("error", err))
		return nil, fmt.Errorf("create room: %w", err)
	}

	return &room, nil
}

// GetRoomByID - получить комнату
func (repo *roomRepository) GetRoomByID(ctx context.Context, id uuid.UUID) (*entity.Room, error) {
	row := repo.db.QueryRow(
		ctx,
		`SELECT id, name, description, capacity, created_at FROM rooms WHERE id = $1`,
		id,
	)

	var room entity.Room

	err := row.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.New(entity.CodeRoomNotFound, entity.ErrMsgRoomNotFound)
		}
		slog.Error("Repository: Error get room",
			slog.Any("roomID", id),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("get room: %w", err)
	}

	return &room, nil
}

// ListRooms - список комнат
func (repo *roomRepository) ListRooms(ctx context.Context) ([]entity.Room, error) {
	rows, err := repo.db.Query(
		ctx,
		`SELECT id, name, description, capacity, created_at FROM rooms ORDER BY created_at`,
	)
	if err != nil {
		slog.Error("Repository: Error list rooms",
			slog.Any("error", err))
		return nil, fmt.Errorf("list rooms: %w", err)
	}
	defer rows.Close()

	var rooms []entity.Room
	for rows.Next() {
		var room entity.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt); err != nil {
			slog.Error("Repository: Error list rooms",
				slog.Any("error", err),
			)
			return nil, fmt.Errorf("list rooms: %w", err)
		}
		rooms = append(rooms, room)
	}
	if rooms == nil {
		rooms = []entity.Room{}
	}

	return rooms, nil
}
