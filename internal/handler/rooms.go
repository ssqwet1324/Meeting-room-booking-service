package handler

import (
	"avito/internal/entity"
	"avito/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListRooms - получить список комнат
func (h *Handler) ListRooms(c *gin.Context) {
	rooms, err := h.rooms.ListRooms(c)
	if err != nil {
		mapError(c, err)
		return
	}

	writeJSON(c, http.StatusOK, gin.H{"rooms": rooms})
}

// CreateRoom - создание комнаты (только админ)
func (h *Handler) CreateRoom(c *gin.Context) {
	role, _ := middleware.GetRole(c)
	if role != entity.RoleAdmin {
		writeError(c, http.StatusForbidden, entity.CodeForbidden, entity.ErrMsgAccessDenied)
		return
	}

	var req entity.NewRoom

	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}
	if req.Name == "" {
		writeError(c, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	room, err := h.rooms.CreateRoom(c, req)
	if err != nil {
		mapError(c, err)
		return
	}

	writeJSON(c, http.StatusCreated, gin.H{"room": room})
}
