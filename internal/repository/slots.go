package repository

import (
	"avito/internal/entity"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type slotsRepository struct {
	db *pgxpool.Pool
}

// NewSlotsRepository создает репозиторий слотов.
func NewSlotsRepository(db *pgxpool.Pool) SlotRepository {
	return &slotsRepository{db: db}
}

// scanSlots - функция для сохранения слотов в список структур
func scanSlots(rows pgx.Rows) ([]entity.Slot, error) {
	var slots []entity.Slot
	for rows.Next() {
		var s entity.Slot
		if err := rows.Scan(&s.ID, &s.RoomID, &s.Start, &s.End); err != nil {
			return nil, fmt.Errorf("scan slot: %w", err)
		}
		slots = append(slots, s)
	}
	if slots == nil {
		slots = []entity.Slot{}
	}
	return slots, rows.Err()
}

// GetSlotByID - получить один слот по Id
func (repo *slotsRepository) GetSlotByID(ctx context.Context, id uuid.UUID) (*entity.Slot, error) {
	rows := repo.db.QueryRow(
		ctx,
		`SELECT id, room_id, start_time, end_time FROM slots WHERE id = $1`,
		id,
	)
	var slot entity.Slot
	err := rows.Scan(&slot.ID, &slot.RoomID, &slot.Start, &slot.End)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.New(entity.CodeSlotNotFound, entity.ErrMsgSlotNotFound)
		}
		slog.Error("Repository: Error getting slot by id",
			slog.Any("id", id),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("get slot by id: %w", err)
	}

	return &slot, nil
}

// GetByRoomAndDate - получить слоты комнаты за день
func (repo *slotsRepository) GetByRoomAndDate(ctx context.Context, roomID uuid.UUID, dayStart, dayEnd time.Time) ([]entity.Slot, error) {
	rows, err := repo.db.Query(
		ctx,
		`SELECT id, room_id, start_time, end_time FROM slots
        WHERE room_id = $1 AND start_time >= $2 AND start_time < $3 
        ORDER BY start_time`,
		roomID, dayStart, dayEnd,
	)
	if err != nil {
		slog.Error("Repository: Error getting slots",
			slog.Any("roomID", roomID),
			slog.Any("dayStart", dayStart),
			slog.Any("dayEnd", dayEnd),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("get slots by date: %w", err)
	}
	defer rows.Close()

	return scanSlots(rows)
}

// GetAvailableByRoomAndDate - получить свободные слоты
func (repo *slotsRepository) GetAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, dayStart, dayEnd time.Time) ([]entity.Slot, error) {
	rows, err := repo.db.Query(
		ctx,
		`SELECT s.id, s.room_id, s.start_time, s.end_time FROM slots s
            WHERE s.room_id = $1 AND s.start_time >= $2 AND s.start_time < $3
            AND NOT EXISTS(
                SELECT 1 FROM bookings b
                WHERE b.slot_id = s.id AND b.status = 'active'
            ) ORDER BY s.start_time`,
		roomID, dayStart, dayEnd,
	)
	if err != nil {
		slog.Error("Repository: Error getting available slots",
			slog.Any("roomID", roomID),
			slog.Any("dayStart", dayStart),
			slog.Any("dayEnd", dayEnd),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("get available slots: %w", err)
	}
	defer rows.Close()

	return scanSlots(rows)
}

// InsertNewSlotsFromRoom - генерация слотов для комнаты на день
func (repo *slotsRepository) InsertNewSlotsFromRoom(ctx context.Context, slots []entity.Slot) error {
	if len(slots) == 0 {
		return nil
	}
	batch := pgx.Batch{}
	for _, slot := range slots {
		batch.Queue(
			`INSERT INTO slots (id, room_id, start_time, end_time) VALUES ($1, $2, $3, $4)
			ON CONFLICT (room_id, start_time) DO NOTHING`,
			slot.ID, slot.RoomID, slot.Start, slot.End,
		)
	}
	br := repo.db.SendBatch(ctx, &batch)

	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {
			slog.Error("Repository: error closing batch", slog.Any("error", err))
		}
	}(br)

	for _, slot := range slots {
		if _, err := br.Exec(); err != nil {
			slog.Error("Repository: Error inserting new slots",
				slog.Any("slot", slot),
				slog.Any("err", err),
			)
			return fmt.Errorf("insert new slots: %w", err)
		}
	}

	return nil
}
