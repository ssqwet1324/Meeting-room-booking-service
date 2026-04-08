package usecase

import (
	"fmt"

	"github.com/google/uuid"
)

// ConferenceUseCase представляет use case конференций.
type ConferenceUseCase struct {
}

// NewConferenceUseCase создает use case конференций.
func NewConferenceUseCase() *ConferenceUseCase {
	return &ConferenceUseCase{}
}

// CreateURL - заглушка для сервиса, который где-то внешне создает ссылку
func (u *ConferenceUseCase) CreateURL(bookingID uuid.UUID) (string, error) {
	return fmt.Sprintf("https://example.com/meeting/%s", bookingID.String()), nil
}
