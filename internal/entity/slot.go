package entity

import (
	"time"

	"github.com/google/uuid"
)

// Slot представляет структуру временного слота.
type Slot struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"roomId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}
