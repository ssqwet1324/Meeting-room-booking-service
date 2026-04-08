package repository

import (
	"avito/internal/config"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// Подключаем file source драйвер для запуска миграций с диска.
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

const (
	maxConnectionsFomPgx  = 20
	minConnectionsFromPgx = 5
)

// Repository представляет структуру доступа к базе данных.
type Repository struct {
	db *pgxpool.Pool
}

// Init - коннект к бд
func Init(ctx context.Context, cfg *config.Config) (*Repository, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.CreateDsn())
	if err != nil {
		slog.Error("Repository: Error get dsn",
			slog.Any("dsn", cfg.CreateDsn()),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("create dsn: %w", err)
	}

	poolCfg.MaxConns = maxConnectionsFomPgx
	poolCfg.MinConns = minConnectionsFromPgx

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		slog.Error("Repository: Error init pool",
			slog.Any("pool", poolCfg),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		slog.Error("Repository: Error ping",
			slog.Any("pool", pool),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("ping pool: %w", err)
	}

	sqlDb := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer func(sqlDb *sql.DB) {
		err := sqlDb.Close()
		if err != nil {
			_ = fmt.Errorf("close db error: %w", err)
		}
	}(sqlDb)

	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		slog.Error("Repository: Error init driver",
			slog.Any("sqlDb", sqlDb),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("init driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		slog.Error("Repository: Error init migrate", slog.Any("error", err))
		return nil, fmt.Errorf("init migrate: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Repository: Error run migration", slog.Any("error", err))
		return nil, fmt.Errorf("run migration: %w", err)
	}

	return &Repository{
		db: pool,
	}, nil
}

// Close закрывает пул подключений к базе данных.
func (repo *Repository) Close() {
	repo.db.Close()
}

// DB возвращает пул подключений pgx.
func (repo *Repository) DB() *pgxpool.Pool {
	return repo.db
}
