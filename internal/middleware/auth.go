package middleware

import (
	"avito/internal/entity"
	"avito/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ContextKeyUserID хранит ключ user_id в контексте
	ContextKeyUserID = "user_id"
	// ContextKeyRole хранит ключ role в контексте
	ContextKeyRole = "role"
)

// AuthMiddleware представляет структуру middleware авторизации.
type AuthMiddleware struct {
	auth *usecase.AuthUseCase
}

// NewAuthMiddleware создает middleware авторизации.
func NewAuthMiddleware(auth *usecase.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

// Authenticate - проверка jwt
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    string(entity.CodeUnauthorized),
					"message": entity.ErrMsgInvalidOrExpiredToken,
				},
			})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims, err := m.auth.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    string(entity.CodeUnauthorized),
					"message": entity.ErrMsgInvalidOrExpiredToken,
				},
			})
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    string(entity.CodeUnauthorized),
					"message": entity.ErrMsgInvalidOrExpiredToken,
				},
			})
			return
		}

		c.Set(ContextKeyUserID, userID)
		c.Set(ContextKeyRole, entity.Role(claims.Role))

		c.Next()
	}
}

// RequireRole - middleware, который требует конкретную роль.
func RequireRole(role entity.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		got, _ := GetRole(c)
		if got != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    string(entity.CodeForbidden),
					"message": entity.ErrMsgAccessDenied,
				},
			})
			return
		}
		c.Next()
	}
}

// GetUserID достаёт userID из gin.Context.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ContextKeyUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

// GetRole достаёт роль из gin.Context.
func GetRole(c *gin.Context) (entity.Role, bool) {
	v, ok := c.Get(ContextKeyRole)
	if !ok {
		return "", false
	}
	r, ok := v.(entity.Role)
	return r, ok
}
