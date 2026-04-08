package usecase

import (
	"avito/internal/entity"
	"avito/internal/repository"
	"context"

	"github.com/google/uuid"
)

// ScheduleUseCase представляет use case расписаний.
type ScheduleUseCase struct {
	schedule repository.ScheduleRepository
	rooms    repository.RoomRepository
}

// NewScheduleUseCase создает use case расписаний.
func NewScheduleUseCase(schedule repository.ScheduleRepository, rooms repository.RoomRepository) *ScheduleUseCase {
	return &ScheduleUseCase{schedule: schedule, rooms: rooms}
}

// GetScheduleByID - получить расписание по id комнаты
func (u *ScheduleUseCase) GetScheduleByID(ctx context.Context, roomID uuid.UUID) (*entity.Schedule, error) {
	return u.schedule.GetScheduleByRoomID(ctx, roomID)
}

// CreateSchedule - создать расписание для комнаты
func (u *ScheduleUseCase) CreateSchedule(ctx context.Context, schedule entity.NewSchedule) (*entity.Schedule, error) {
	for _, d := range schedule.DaysOfWeek {
		if d < 1 || d > 7 {
			return nil, entity.New(entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		}
	}
	if len(schedule.DaysOfWeek) == 0 {
		return nil, entity.New(entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
	}

	if _, err := u.rooms.GetRoomByID(ctx, schedule.RoomID); err != nil {
		return nil, err
	}

	newSchedule, err := u.schedule.CreateSchedule(ctx, schedule)
	if err != nil {
		return nil, err
	}

	return newSchedule, nil
}
