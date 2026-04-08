package entity

// Тексты ошибок API
const (
	ErrMsgInvalidRequest        = "invalid request"
	ErrMsgInternalServerError   = "internal server error"
	ErrMsgInvalidCredentials    = "invalid credentials"
	ErrMsgInvalidOrExpiredToken = "invalid or expired token"
	ErrMsgAccessDenied          = "access denied"
	ErrMsgScheduleExists        = "schedule for this room already exists and cannot be changed"
	ErrMsgSlotAlreadyBooked     = "slot is already booked"
	ErrMsgCancelOtherUser       = "cannot cancel another user's booking"
	ErrMsgRoomNotFound          = "room not found"
	ErrMsgSlotNotFound          = "slot not found"
	ErrMsgBookingNotFound       = "booking not found"
	ErrMsgNotFound              = "not found"
)
