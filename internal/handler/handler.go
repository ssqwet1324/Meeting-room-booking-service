package handler

import (
	"avito/internal/entity"
	"avito/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler представляет структуру HTTP-обработчиков.
type Handler struct {
	auth     *usecase.AuthUseCase
	bookings *usecase.BookingUseCase
	rooms    *usecase.RoomUseCase
	schedule *usecase.ScheduleUseCase
	slots    *usecase.SlotUseCase
}

// New создает новый обработчик.
func New(
	auth *usecase.AuthUseCase,
	rooms *usecase.RoomUseCase,
	schedule *usecase.ScheduleUseCase,
	slots *usecase.SlotUseCase,
	bookings *usecase.BookingUseCase,
) *Handler {
	return &Handler{auth: auth,
		rooms:    rooms,
		schedule: schedule,
		slots:    slots,
		bookings: bookings,
	}
}

// mapError - замапить ошибку и отправить в ответ
func mapError(c *gin.Context, err error) {
	code, ok := entity.GetCode(err)
	if !ok {
		writeError(c, http.StatusInternalServerError, entity.CodeInternalError, entity.ErrMsgInternalServerError)
		return
	}

	status := codeToStatus(code)
	writeError(c, status, code, err.Error())
}

// writeError - формируем json для ошибки
func writeError(c *gin.Context, status int, code entity.Code, msg string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    string(code),
			"message": msg,
		},
	})
	c.Abort()
}

// writeJSON универсальная обертка для json
func writeJSON(c *gin.Context, status int, v any) {
	c.JSON(status, v)
}

// codeToStatus - конвертация ошибки в HTTP ошибку
func codeToStatus(code entity.Code) int {
	switch code {
	case entity.CodeInvalidRequest:
		return http.StatusBadRequest
	case entity.CodeUnauthorized:
		return http.StatusUnauthorized
	case entity.CodeForbidden:
		return http.StatusForbidden
	case entity.CodeNotFound, entity.CodeRoomNotFound, entity.CodeSlotNotFound, entity.CodeBookingNotFound:
		return http.StatusNotFound
	case entity.CodeSlotBooked, entity.CodeScheduleExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
