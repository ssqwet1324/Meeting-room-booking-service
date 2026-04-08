package entity

import (
	"time"

	"github.com/google/uuid"
)

// Room представляет структуру комнаты
type Room struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Capacity    *int       `json:"capacity"`
	CreatedAt   *time.Time `json:"createdAt"`
}

// Schedule представляет структуру расписания комнаты
type Schedule struct {
	ID         uuid.UUID `json:"id"`
	RoomID     uuid.UUID `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek"`
	StartTime  string    `json:"startTime"`
	EndTime    string    `json:"endTime"`
}

// NewRoom представляет структуру запроса на создание комнаты
type NewRoom struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

// NewSchedule представляет структуру запроса на создание расписания
type NewSchedule struct {
	RoomID     uuid.UUID `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek"`
	StartTime  string    `json:"startTime"`
	EndTime    string    `json:"endTime"`
}
