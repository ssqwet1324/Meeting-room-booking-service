package usecase

import (
	"avito/internal/entity"
	"avito/internal/repository"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// AdminUUID задает фиксированный UUID администратора для dummy-логина
	AdminUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	// UserUUID задает фиксированный UUID пользователя для dummy-логина
	UserUUID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

// AuthUseCase представляет use case авторизации.
type AuthUseCase struct {
	users     repository.UserRepository
	jwtSecret string
}

// NewAuthUseCase создает use case авторизации.
func NewAuthUseCase(users repository.UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		users:     users,
		jwtSecret: jwtSecret,
	}
}

// DummyLogin - логика для тестовой отдачи jwt
func (u *AuthUseCase) DummyLogin(role entity.Role) (string, error) {
	var userID uuid.UUID
	var email string
	switch role {
	case entity.RoleAdmin:
		userID = AdminUUID
		email = "admin@example.com"
	case entity.RoleUser:
		userID = UserUUID
		email = "user@example.com"
	default:
		return "", entity.New(entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
	}
	_ = email

	return u.issueToken(userID, role)
}

// Register - регистрация пользователя
func (u *AuthUseCase) Register(ctx context.Context, email, password string, role entity.Role) (*entity.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error hashing password", slog.Any("error", err))
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newUser, err := u.users.Create(ctx, email, string(hashPassword), role)
	if err != nil {
		slog.Error("Error creating user", slog.Any("error", err))
		return nil, err
	}

	return newUser, nil
}

// Login - логин
func (u *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		return "", entity.New(entity.CodeUnauthorized, entity.ErrMsgInvalidCredentials)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", entity.New(entity.CodeUnauthorized, entity.ErrMsgInvalidCredentials)
	}

	return u.issueToken(user.ID, user.Role)
}

// issueToken - выдаем пользователю новый jwt токен
func (u *AuthUseCase) issueToken(userID uuid.UUID, role entity.Role) (string, error) {
	claims := &entity.Claims{
		UserID: userID.String(),
		Role:   string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		slog.Error("Error signing token",
			slog.Any("userID", userID),
			slog.Any("error", err),
		)
		return "", fmt.Errorf("signing token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken - валидация токена
func (u *AuthUseCase) ValidateToken(tokenString string) (*entity.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(u.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, entity.New(entity.CodeUnauthorized, entity.ErrMsgInvalidOrExpiredToken)
	}
	claims, ok := token.Claims.(*entity.Claims)
	if !ok {
		slog.Error("Error parsing token claims", slog.Any("error", err))
		return nil, entity.New(entity.CodeUnauthorized, entity.ErrMsgInvalidOrExpiredToken)
	}

	return claims, nil
}
