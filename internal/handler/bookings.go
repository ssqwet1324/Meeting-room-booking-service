package handler

import (
	"avito/internal/entity"
	"avito/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBooking - создание брони
func (h *Handler) CreateBooking(ctx *gin.Context) {
	role, _ := middleware.GetRole(ctx)
	if role != entity.RoleUser {
		writeError(ctx, http.StatusForbidden, entity.CodeForbidden, entity.ErrMsgAccessDenied)
		return
	}

	userID, _ := middleware.GetUserID(ctx)

	var req entity.CreateBookingRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	slotID, err := uuid.Parse(req.SlotID)
	if err != nil {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	booking, err := h.bookings.Create(ctx, userID, slotID, req.CreateConferenceLink)
	if err != nil {
		mapError(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusCreated, gin.H{"booking": booking})
}

// ListBookings - список всех броней (для админа)
func (h *Handler) ListBookings(ctx *gin.Context) {
	role, _ := middleware.GetRole(ctx)
	if role != entity.RoleAdmin {
		writeError(ctx, http.StatusForbidden, entity.CodeForbidden, entity.ErrMsgAccessDenied)
		return
	}

	page := 1
	pageSize := 20

	if v := ctx.Query("page"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil || p < 1 {
			writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
			return
		}
		page = p
	}

	if v := ctx.Query("pageSize"); v != "" {
		ps, err := strconv.Atoi(v)
		if err != nil || ps < 1 || ps > 100 {
			writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
			return
		}
		pageSize = ps
	}

	bookings, total, err := h.bookings.ListAll(ctx, page, pageSize)
	if err != nil {
		mapError(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusOK, gin.H{
		"bookings": bookings,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}

// MyBookings - брони конкретного пользователя
func (h *Handler) MyBookings(ctx *gin.Context) {
	role, _ := middleware.GetRole(ctx)
	if role != entity.RoleUser {
		writeError(ctx, http.StatusForbidden, entity.CodeForbidden, entity.ErrMsgAccessDenied)
		return
	}

	userID, _ := middleware.GetUserID(ctx)

	bookings, err := h.bookings.MyBookings(ctx, userID)
	if err != nil {
		mapError(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusOK, gin.H{"bookings": bookings})
}

// CancelBooking - отмена брони
func (h *Handler) CancelBooking(ctx *gin.Context) {
	role, _ := middleware.GetRole(ctx)
	if role != entity.RoleUser {
		writeError(ctx, http.StatusForbidden, entity.CodeForbidden, entity.ErrMsgAccessDenied)
		return
	}

	bookingIDStr := ctx.Param("bookingId")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	userID, _ := middleware.GetUserID(ctx)

	booking, err := h.bookings.Cancel(ctx, bookingID, userID)
	if err != nil {
		mapError(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusOK, gin.H{"booking": booking})
}
