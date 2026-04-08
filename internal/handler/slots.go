package handler

import (
	"avito/internal/entity"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListSlots - доступные слоты для комнаты на дату
func (h *Handler) ListSlots(c *gin.Context) {
	roomIDStr := c.Param("roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		writeError(c, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	slots, err := h.slots.GetAvailableSlots(c, roomID, date.UTC())
	if err != nil {
		mapError(c, err)
		return
	}

	writeJSON(c, http.StatusOK, gin.H{"slots": slots})
}
