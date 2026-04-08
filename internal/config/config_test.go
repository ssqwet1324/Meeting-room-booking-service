package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_CreateDsn(t *testing.T) {
	cfg := &Config{
		DbUser:     "u",
		DbPassword: "p",
		DbHost:     "localhost",
		DbPort:     5432,
		DbName:     "db",
	}
	dsn := cfg.CreateDsn()
	assert.Contains(t, dsn, "postgres://u:p@localhost:5432/db")
	assert.Contains(t, dsn, "sslmode=disable")
}

func TestNew_ReadsEnv(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_USER", "user")
	t.Setenv("DB_PASSWORD", "pass")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "5432")

	cfg, err := New()
	require.NoError(t, err)
	assert.Equal(t, "test-secret", cfg.JWTSecret)
	assert.Equal(t, "testdb", cfg.DbName)
	assert.Equal(t, 5432, cfg.DbPort)
}
