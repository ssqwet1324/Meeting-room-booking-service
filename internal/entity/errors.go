package entity

import "errors"

// Code описывает код прикладной ошибки
type Code string

const (
	// CodeInvalidRequest задает код ошибки невалидного запроса
	CodeInvalidRequest Code = "INVALID_REQUEST"
	// CodeUnauthorized задает код ошибки авторизации
	CodeUnauthorized Code = "UNAUTHORIZED"
	// CodeNotFound задает код ошибки отсутствия сущности
	CodeNotFound Code = "NOT_FOUND"
	// CodeRoomNotFound задает код ошибки отсутствия комнаты
	CodeRoomNotFound Code = "ROOM_NOT_FOUND"
	// CodeSlotNotFound задает код ошибки отсутствия слота
	CodeSlotNotFound Code = "SLOT_NOT_FOUND"
	// CodeSlotBooked задает код ошибки занятого слота
	CodeSlotBooked Code = "SLOT_ALREADY_BOOKED"
	// CodeBookingNotFound задает код ошибки отсутствия брони
	CodeBookingNotFound Code = "BOOKING_NOT_FOUND"
	// CodeForbidden задает код ошибки запрета доступа
	CodeForbidden Code = "FORBIDDEN"
	// CodeScheduleExists задает код ошибки существующего расписания
	CodeScheduleExists Code = "SCHEDULE_EXISTS"
	// CodeInternalError задает код внутренней ошибки
	CodeInternalError Code = "INTERNAL_ERROR"
)

// AppError представляет структуру прикладной ошибки
type AppError struct {
	Code    Code
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// New создает новую прикладную ошибку.
func New(code Code, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

// GetCode извлекает код прикладной ошибки из error.
func GetCode(err error) (Code, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code, true
	}
	return "", false
}
