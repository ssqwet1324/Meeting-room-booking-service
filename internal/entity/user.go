package entity

import (
	"time"

	"github.com/google/uuid"
)

// Role описывает роль пользователя
type Role string

const (
	// RoleAdmin задает роль администратора
	RoleAdmin Role = "admin"
	// RoleUser задает роль обычного пользователя
	RoleUser Role = "user"
)

// Roles представляет структуру роли в ответе API
type Roles struct {
	Role string `json:"role"`
}

// User представляет структуру пользователя
type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Role         Role       `json:"role"`
	PasswordHash string     `json:"-"`
	CreatedAt    *time.Time `json:"createdAt"`
}

// Register представляет структуру запроса на регистрацию
type Register struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Login представляет структуру запроса на вход
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
