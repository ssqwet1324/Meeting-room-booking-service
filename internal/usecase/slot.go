package usecase

import (
	"avito/internal/entity"
	"avito/internal/repository"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// slotDuration - длительность слота
const slotDuration = 30 * time.Minute

// SlotUseCase представляет use case слотов.
type SlotUseCase struct {
	slots     repository.SlotRepository
	schedules repository.ScheduleRepository
	rooms     repository.RoomRepository
}

// NewSlotUseCase создает use case слотов.
func NewSlotUseCase(slots repository.SlotRepository,
	schedules repository.ScheduleRepository,
	rooms repository.RoomRepository) *SlotUseCase {
	return &SlotUseCase{
		slots:     slots,
		schedules: schedules,
		rooms:     rooms,
	}
}

// GetAvailableSlots - получить свободные слоты
func (u *SlotUseCase) GetAvailableSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]entity.Slot, error) {
	// проверяем, что комната существует
	if _, err := u.rooms.GetRoomByID(ctx, roomID); err != nil {
		return nil, err
	}

	// проверяем есть ли расписание, если нет, значит мест свободных нет
	schedule, err := u.schedules.GetScheduleByRoomID(ctx, roomID)
	if err != nil {
		if code, ok := entity.GetCode(err); ok && code == entity.CodeNotFound {
			return []entity.Slot{}, nil
		}
		return nil, err
	}

	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	weekday := toAPIWeekday(dayStart.Weekday())
	if !containsDay(schedule.DaysOfWeek, weekday) {
		return []entity.Slot{}, nil
	}

	exist, err := u.slots.GetByRoomAndDate(ctx, roomID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	if len(exist) == 0 {
		generated := generateSlots(roomID, dayStart, schedule)
		if err := u.slots.InsertNewSlotsFromRoom(ctx, generated); err != nil {
			return nil, err
		}
	}

	availableSlots, err := u.slots.GetAvailableByRoomAndDate(ctx, roomID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	return availableSlots, nil
}

// GetSlotByID - получаем слот по id
func (u *SlotUseCase) GetSlotByID(ctx context.Context, id uuid.UUID) (*entity.Slot, error) {
	return u.slots.GetSlotByID(ctx, id)
}

// generateSlots - генерирует 30-минутные слоты по расписанию
func generateSlots(roomID uuid.UUID, dayStart time.Time, schedule *entity.Schedule) []entity.Slot {
	startH, startM := parseTime(schedule.StartTime)
	endH, endM := parseTime(schedule.EndTime)

	slotStart := dayStart.Add(time.Duration(startH)*time.Hour + time.Duration(startM)*time.Minute)
	schedEnd := dayStart.Add(time.Duration(endH)*time.Hour + time.Duration(endM)*time.Minute)

	var slots []entity.Slot
	for t := slotStart; !t.Add(slotDuration).After(schedEnd); t = t.Add(slotDuration) {
		slots = append(slots, entity.Slot{
			ID:     uuid.New(),
			RoomID: roomID,
			Start:  t,
			End:    t.Add(slotDuration),
		})
	}

	return slots
}

func parseTime(hhmm string) (int, int) {
	parts := strings.SplitN(hhmm, ":", 2)
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return h, m
}

// toAPIWeekday converts Go's time.Weekday to API weekday (1=Mon..7=Sun).
func toAPIWeekday(wd time.Weekday) int {
	if wd == time.Sunday {
		return 7
	}
	return int(wd)
}

func containsDay(days []int, day int) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}
