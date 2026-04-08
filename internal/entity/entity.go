package entity

import (
	"github.com/golang-jwt/jwt/v4"
)

// Claims представляет JWT claims пользователя.
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
