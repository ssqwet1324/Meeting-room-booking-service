package config

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - структура конфига
type Config struct {
	JWTSecret  string `env:"JWT_SECRET"`
	DbName     string `env:"DB_NAME"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbHost     string `env:"DB_HOST"`
	DbPort     int    `env:"DB_PORT"`
}

// New - конструктор
func New() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		slog.Error("Config: Error reading env file: ",
			slog.Any("cfg", cfg),
			slog.Any("error", err),
		)
		return nil, err
	}

	return &cfg, nil
}

// CreateDsn - создание строки подключения
func (cfg *Config) CreateDsn() string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		strconv.Itoa(cfg.DbPort),
		cfg.DbName,
	)

	return dsn
}
