package repository

import (
	"avito/internal/entity"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type scheduleRepository struct {
	db *pgxpool.Pool
}

// NewScheduleRepository создает репозиторий расписаний.
func NewScheduleRepository(db *pgxpool.Pool) ScheduleRepository {
	return &scheduleRepository{db: db}
}

// CreateSchedule - создать расписание для комнаты
func (repo *scheduleRepository) CreateSchedule(ctx context.Context, schedule entity.NewSchedule) (*entity.Schedule, error) {
	row := repo.db.QueryRow(ctx,
		`INSERT INTO schedules (room_id, days_of_week, start_time, end_time)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, room_id, days_of_week, start_time, end_time`,
		schedule.RoomID, schedule.DaysOfWeek, schedule.StartTime, schedule.EndTime)

	var s entity.Schedule
	var retDays []int

	if err := row.Scan(&s.ID, &s.RoomID, &retDays, &s.StartTime, &s.EndTime); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entity.New(entity.CodeScheduleExists, entity.ErrMsgScheduleExists)
		}
		slog.Error("CreateSchedule: scan failed",
			slog.Any("schedule", schedule),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("create schedule: %w", err)
	}

	s.DaysOfWeek = retDays

	return &s, nil
}

// GetScheduleByRoomID возвращает расписание для комнаты по её roomID.
func (repo *scheduleRepository) GetScheduleByRoomID(ctx context.Context, roomID uuid.UUID) (*entity.Schedule, error) {
	row := repo.db.QueryRow(
		ctx,
		`SELECT id, room_id, days_of_week, start_time, end_time FROM schedules WHERE room_id = $1`,
		roomID,
	)

	var s entity.Schedule
	var days []int

	if err := row.Scan(&s.ID, &s.RoomID, &days, &s.StartTime, &s.EndTime); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.New(entity.CodeNotFound, entity.ErrMsgNotFound)
		}
		slog.Error("GetScheduleByID: scan failed",
			slog.Any("roomID", roomID),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("get schedule: %w", err)
	}

	s.DaysOfWeek = days

	return &s, nil
}
