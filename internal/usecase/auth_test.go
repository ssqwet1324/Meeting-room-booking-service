package usecase

import (
	"context"
	"testing"

	"avito/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyLogin_AdminReturnsToken(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{}, "secret")
	token, err := uc.DummyLogin(entity.RoleAdmin)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestDummyLogin_UserReturnsToken(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{}, "secret")
	token, err := uc.DummyLogin(entity.RoleUser)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{}, "secret")
	_, err := uc.DummyLogin(entity.Role("superuser"))
	require.Error(t, err)
	code, ok := entity.GetCode(err)
	require.True(t, ok)
	assert.Equal(t, entity.CodeInvalidRequest, code)
}

func TestDummyLogin_AdminFixedUUID(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{}, "testsecret")
	token, err := uc.DummyLogin(entity.RoleAdmin)
	require.NoError(t, err)
	claims, err := uc.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, AdminUUID.String(), claims.UserID)
	assert.Equal(t, string(entity.RoleAdmin), claims.Role)
}

func TestDummyLogin_UserFixedUUID(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{}, "testsecret")
	token, err := uc.DummyLogin(entity.RoleUser)
	require.NoError(t, err)
	claims, err := uc.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, UserUUID.String(), claims.UserID)
	assert.Equal(t, string(entity.RoleUser), claims.Role)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{}, "secret")
	_, err := uc.ValidateToken("not.a.token")
	require.Error(t, err)
	code, ok := entity.GetCode(err)
	require.True(t, ok)
	assert.Equal(t, entity.CodeUnauthorized, code)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	uc1 := NewAuthUseCase(&mockUserRepo{}, "secret1")
	uc2 := NewAuthUseCase(&mockUserRepo{}, "secret2")
	token, err := uc1.DummyLogin(entity.RoleUser)
	require.NoError(t, err)
	_, err = uc2.ValidateToken(token)
	require.Error(t, err)
}

func TestRegister_Success(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{}, "secret")
	user, err := uc.Register(context.Background(), "test@example.com", "password123", entity.RoleUser)
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	uc := NewAuthUseCase(&mockUserRepo{err: entity.New(entity.CodeNotFound, "not found")}, "secret")
	_, err := uc.Login(context.Background(), "x@y.com", "pw")
	require.Error(t, err)
	code, ok := entity.GetCode(err)
	require.True(t, ok)
	assert.Equal(t, entity.CodeUnauthorized, code)
}
