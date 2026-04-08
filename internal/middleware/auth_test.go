package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito/internal/entity"
	"avito/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubUserRepo struct {
	user *entity.User
	err  error
}

func (s *stubUserRepo) GetByID(_ context.Context, _ uuid.UUID) (*entity.User, error) {
	return s.user, s.err
}

func (s *stubUserRepo) GetByEmail(_ context.Context, _ string) (*entity.User, error) {
	return s.user, s.err
}

func (s *stubUserRepo) Create(_ context.Context, _, _ string, _ entity.Role) (*entity.User, error) {
	return nil, nil
}

func TestAuthenticate_MissingBearer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "secret")
	mw := NewAuthMiddleware(auth)

	r := gin.New()
	r.Use(mw.Authenticate())
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "secret")
	mw := NewAuthMiddleware(auth)

	r := gin.New()
	r.Use(mw.Authenticate())
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_Success_SetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "test-secret")
	token, err := auth.DummyLogin(entity.RoleUser)
	require.NoError(t, err)

	mw := NewAuthMiddleware(auth)

	r := gin.New()
	r.Use(mw.Authenticate())
	r.GET("/", func(c *gin.Context) {
		uid, ok := GetUserID(c)
		require.True(t, ok)
		assert.Equal(t, usecase.UserUUID, uid)
		role, ok := GetRole(c)
		require.True(t, ok)
		assert.Equal(t, entity.RoleUser, role)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_Allows(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.Set(ContextKeyRole, entity.RoleAdmin)
		c.Next()
	}, RequireRole(entity.RoleAdmin), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRole_Denies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.Set(ContextKeyRole, entity.RoleUser)
		c.Next()
	}, RequireRole(entity.RoleAdmin), func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
