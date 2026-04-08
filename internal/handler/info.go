package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Info - ручка для проверки сервиса
func (h *Handler) Info(ctx *gin.Context) {
	writeJSON(ctx, http.StatusOK, gin.H{"status": "ok"})
}
