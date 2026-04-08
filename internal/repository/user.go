package repository

import (
	"avito/internal/entity"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает репозиторий пользователей.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// GetByID - получить пользователя по id
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row := r.db.QueryRow(
		ctx,
		`SELECT id, email, role, password_hash, created_at FROM users WHERE id = $1`,
		id,
	)

	return scanUser(row)
}

// GetByEmail - получить пользователя по почте
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, role, password_hash, created_at FROM users WHERE email = $1`, email)

	return scanUser(row)
}

// Create - регистрация пользователя
func (r *userRepository) Create(ctx context.Context, email, passwordHash string, role entity.Role) (*entity.User, error) {
	row := r.db.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)
		 RETURNING id, email, role, password_hash, created_at`,
		email, passwordHash, role)

	u, err := scanUser(row)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entity.New(entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest)
		}
		slog.Error("Error scanning user:",
			slog.String("email", email),
			slog.Any("error", err))
		return nil, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}

// scanUser - сканирование полей с бд в структуру
func scanUser(row pgx.Row) (*entity.User, error) {
	var u entity.User

	err := row.Scan(&u.ID, &u.Email, &u.Role, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.New(entity.CodeNotFound, entity.ErrMsgNotFound)
		}
		return nil, fmt.Errorf("scan user: %w", err)
	}

	return &u, nil
}
