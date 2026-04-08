package entity

import (
	"time"

	"github.com/google/uuid"
)

// BookingStatus описывает статус бронирования
type BookingStatus string

const (
	// BookingStatusActive задает активный статус брони
	BookingStatusActive BookingStatus = "active"
	// BookingStatusCancelled задает отмененный статус брони
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Booking представляет структуру бронирования
type Booking struct {
	ID             uuid.UUID     `json:"id"`
	SlotID         uuid.UUID     `json:"slotId"`
	UserID         uuid.UUID     `json:"userId"`
	Status         BookingStatus `json:"status"`
	ConferenceLink *string       `json:"conferenceLink"`
	CreatedAt      *time.Time    `json:"createdAt"`
}

// CreateBookingRequest представляет запрос на создание брони
type CreateBookingRequest struct {
	SlotID               string `json:"slotId"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}
