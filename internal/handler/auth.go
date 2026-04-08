package handler

import (
	"avito/internal/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DummyLogin - ручка получения тестового jwt
func (h *Handler) DummyLogin(ctx *gin.Context) {
	var req entity.Roles

	if err := ctx.ShouldBind(&req); err != nil {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	role := entity.Role(req.Role)
	if role != entity.RoleAdmin && role != entity.RoleUser {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	token, err := h.auth.DummyLogin(role)
	if err != nil {
		mapError(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusOK, gin.H{"token": token})
}

// Register - ручка регистрации
func (h *Handler) Register(ctx *gin.Context) {
	var req entity.Register

	if err := ctx.ShouldBind(&req); err != nil {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	role := entity.Role(req.Role)
	if role != entity.RoleAdmin && role != entity.RoleUser {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	user, err := h.auth.Register(ctx, req.Email, req.Password, role)
	if err != nil {
		mapError(ctx, err)
		return
	}

	fmt.Println("user:", user)

	writeJSON(ctx, http.StatusCreated, gin.H{"user": user})
}

// Login - Ручка для входа
func (h *Handler) Login(ctx *gin.Context) {
	var req entity.Login

	if err := ctx.ShouldBind(&req); err != nil {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		writeError(ctx, http.StatusBadRequest, entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		return
	}

	token, err := h.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		mapError(ctx, err)
		return
	}

	writeJSON(ctx, http.StatusOK, gin.H{"token": token})
}
