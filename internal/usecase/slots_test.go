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

func TestGetAvailableSlots_NoSchedule_ReturnsEmpty(t *testing.T) {
	roomID := uuid.New()
	uc := NewSlotUseCase(
		&mockSlotRepoSlots{},
		&mockScheduleRepo{err: entity.New(entity.CodeNotFound, "schedule not found")},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)

	slots, err := uc.GetAvailableSlots(context.Background(), roomID, time.Now())
	require.NoError(t, err)
	assert.Empty(t, slots)
}

func TestGetAvailableSlots_WrongWeekday_ReturnsEmpty(t *testing.T) {
	roomID := uuid.New()
	schedule := &entity.Schedule{DaysOfWeek: []int{1}, StartTime: "09:00", EndTime: "18:00"}
	sunday := time.Date(2024, 6, 9, 0, 0, 0, 0, time.UTC)

	uc := NewSlotUseCase(
		&mockSlotRepoSlots{},
		&mockScheduleRepo{schedule: schedule},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)

	slots, err := uc.GetAvailableSlots(context.Background(), roomID, sunday)
	require.NoError(t, err)
	assert.Empty(t, slots)
}

func TestGetAvailableSlots_RoomNotFound_ReturnsError(t *testing.T) {
	uc := NewSlotUseCase(
		&mockSlotRepoSlots{},
		&mockScheduleRepo{},
		&mockRoomRepo{err: entity.New(entity.CodeRoomNotFound, "room not found")},
	)

	_, err := uc.GetAvailableSlots(context.Background(), uuid.New(), time.Now())
	require.Error(t, err)
	code, _ := entity.GetCode(err)
	assert.Equal(t, entity.CodeRoomNotFound, code)
}

func TestGetAvailableSlots_GeneratesSlots_WhenNoneExist(t *testing.T) {
	roomID := uuid.New()
	schedule := &entity.Schedule{DaysOfWeek: []int{1, 2, 3, 4, 5}, StartTime: "09:00", EndTime: "11:00"}
	monday := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)

	slotRepo := &mockSlotRepoSlots{existing: []entity.Slot{}, available: []entity.Slot{}}

	uc := NewSlotUseCase(
		slotRepo,
		&mockScheduleRepo{schedule: schedule},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)

	_, err := uc.GetAvailableSlots(context.Background(), roomID, monday)
	require.NoError(t, err)
	// 09:00–11:00 = 4 слота по 30 минут
	assert.Len(t, slotRepo.inserted, 4)
}

func TestGetAvailableSlots_DoesNotRegenerateIfSlotsExist(t *testing.T) {
	roomID := uuid.New()
	schedule := &entity.Schedule{DaysOfWeek: []int{1}, StartTime: "09:00", EndTime: "10:00"}
	monday := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)

	existingSlot := entity.Slot{ID: uuid.New(), RoomID: roomID}
	slotRepo := &mockSlotRepoSlots{existing: []entity.Slot{existingSlot}}

	uc := NewSlotUseCase(
		slotRepo,
		&mockScheduleRepo{schedule: schedule},
		&mockRoomRepo{room: &entity.Room{ID: roomID}},
	)

	_, err := uc.GetAvailableSlots(context.Background(), roomID, monday)
	require.NoError(t, err)
	assert.Nil(t, slotRepo.inserted)
}

func TestGenerateSlots_CorrectCount(t *testing.T) {
	roomID := uuid.New()
	schedule := &entity.Schedule{StartTime: "09:00", EndTime: "18:00"}
	day := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	slots := generateSlots(roomID, day, schedule)
	assert.Len(t, slots, 18)
	assert.Equal(t, time.Hour*9, slots[0].Start.Sub(day))
	assert.Equal(t, 30*time.Minute, slots[0].End.Sub(slots[0].Start))
}

func TestToAPIWeekday(t *testing.T) {
	assert.Equal(t, 1, toAPIWeekday(time.Monday))
	assert.Equal(t, 7, toAPIWeekday(time.Sunday))
	assert.Equal(t, 6, toAPIWeekday(time.Saturday))
}

func TestContainsDay(t *testing.T) {
	assert.True(t, containsDay([]int{1, 3, 5}, 3))
	assert.False(t, containsDay([]int{1, 2}, 7))
}

func TestSlotUseCase_GetSlotByID(t *testing.T) {
	id := uuid.New()
	slot := &entity.Slot{ID: id, RoomID: uuid.New()}
	repo := &mockSlotRepoSlots{slotByID: slot}
	uc := NewSlotUseCase(repo, &mockScheduleRepo{}, &mockRoomRepo{})

	got, err := uc.GetSlotByID(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, slot, got)
}
