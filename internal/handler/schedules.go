package handler

import (
	"avito/internal/entity"
	"avito/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateSchedule - создание расписания комнаты (только админ)
func (h *Handler) CreateSchedule(c *gin.Context) {
	role, _ := middleware.GetRole(c)
	if role != entity.RoleAdmin {
		writeError(c, http.StatusForbidden, entity.CodeForbidden, entity.ErrMsgAccessDenied)
		return
	}

	roomIDStr := c.Param("roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	var req entity.NewSchedule
	req.RoomID = roomID

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	schedule, err := h.schedule.CreateSchedule(c, req)
	if err != nil {
		mapError(c, err)
		return
	}

	writeJSON(c, http.StatusCreated, gin.H{"schedule": schedule})
}
