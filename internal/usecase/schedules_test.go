package usecase

import (
	"context"
	"testing"

	"avito/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSchedule_InvalidDayOfWeek(t *testing.T) {
	roomID := uuid.New()
	uc := NewScheduleUseCase(
		&mockScheduleRepo{},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)
	_, err := uc.CreateSchedule(context.Background(), entity.NewSchedule{
		RoomID: roomID, DaysOfWeek: []int{0, 1, 2}, StartTime: "09:00", EndTime: "18:00",
	})
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeInvalidRequest, code)
}

func TestCreateSchedule_DayOutOfRange(t *testing.T) {
	roomID := uuid.New()
	uc := NewScheduleUseCase(
		&mockScheduleRepo{},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)
	_, err := uc.CreateSchedule(context.Background(), entity.NewSchedule{
		RoomID: roomID, DaysOfWeek: []int{1, 8}, StartTime: "09:00", EndTime: "18:00",
	})
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeInvalidRequest, code)
}

func TestCreateSchedule_EmptyDays(t *testing.T) {
	roomID := uuid.New()
	uc := NewScheduleUseCase(
		&mockScheduleRepo{},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)
	_, err := uc.CreateSchedule(context.Background(), entity.NewSchedule{
		RoomID: roomID, DaysOfWeek: []int{}, StartTime: "09:00", EndTime: "18:00",
	})
	require.Error(t, err)
}

func TestCreateSchedule_RoomNotFound(t *testing.T) {
	uc := NewScheduleUseCase(
		&mockScheduleRepo{},
		&mockRoomRepo{err: entity.New(entity.CodeRoomNotFound, "room not found")},
	)
	_, err := uc.CreateSchedule(context.Background(), entity.NewSchedule{
		RoomID: uuid.New(), DaysOfWeek: []int{1}, StartTime: "09:00", EndTime: "18:00",
	})
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeRoomNotFound, code)
}

func TestCreateSchedule_AlreadyExists(t *testing.T) {
	roomID := uuid.New()
	uc := NewScheduleUseCase(
		&mockScheduleRepo{err: entity.New(entity.CodeScheduleExists, "schedule exists")},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)
	_, err := uc.CreateSchedule(context.Background(), entity.NewSchedule{
		RoomID: roomID, DaysOfWeek: []int{1, 2, 3}, StartTime: "09:00", EndTime: "18:00",
	})
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeScheduleExists, code)
}

func TestCreateSchedule_Success(t *testing.T) {
	roomID := uuid.New()
	uc := NewScheduleUseCase(
		&mockScheduleRepo{},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)
	s, err := uc.CreateSchedule(context.Background(), entity.NewSchedule{
		RoomID: roomID, DaysOfWeek: []int{1, 2, 3, 4, 5}, StartTime: "09:00", EndTime: "18:00",
	})
	require.NoError(t, err)
	assert.Equal(t, roomID, s.RoomID)
}

func TestGetScheduleByID(t *testing.T) {
	roomID := uuid.New()
	sch := &entity.Schedule{RoomID: roomID, DaysOfWeek: []int{1}}
	uc := NewScheduleUseCase(
		&mockScheduleRepo{schedule: sch},
		&mockRoomRepo{},
	)
	got, err := uc.GetScheduleByID(context.Background(), roomID)
	require.NoError(t, err)
	assert.Equal(t, roomID, got.RoomID)
}
