package repository

import (
	"avito/internal/entity"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bookingRepository struct {
	db *pgxpool.Pool
}

// NewBookingRepository создает репозиторий бронирований.
func NewBookingRepository(db *pgxpool.Pool) BookingRepository {
	return &bookingRepository{db: db}
}

// scanBookings - функция для сохранения броней в список структур
func scanBookings(rows pgx.Rows) ([]entity.Booking, error) {
	var bookings []entity.Booking
	for rows.Next() {
		var b entity.Booking
		if err := rows.Scan(&b.ID, &b.SlotID, &b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		bookings = append(bookings, b)
	}
	if bookings == nil {
		bookings = []entity.Booking{}
	}

	return bookings, rows.Err()
}

// scanBooking функция для возврата готовой структуры
func scanBooking(row pgx.Row) (*entity.Booking, error) {
	var b entity.Booking
	err := row.Scan(&b.ID, &b.SlotID, &b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.New(entity.CodeBookingNotFound, entity.ErrMsgBookingNotFound)
		}
		return nil, fmt.Errorf("scan booking: %w", err)
	}
	return &b, nil
}

// CreateBooking - создать бронь
func (repo *bookingRepository) CreateBooking(ctx context.Context, booking *entity.Booking) (*entity.Booking, error) {
	rows := repo.db.QueryRow(
		ctx,
		`INSERT INTO bookings (slot_id, user_id, status, conference_link) VALUES ($1, $2, $3, $4)
 			RETURNING id, slot_id, user_id, status, conference_link, created_at`,
		booking.SlotID, booking.UserID, booking.Status, booking.ConferenceLink,
	)

	b, err := scanBooking(rows)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entity.New(entity.CodeSlotBooked, entity.ErrMsgSlotAlreadyBooked)
		}
		slog.Error("Error inserting booking",
			slog.Any("booking", booking),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("error scan booking: %w", err)
	}

	return b, nil
}

// GetBookingByID - получить бронь по id
func (repo *bookingRepository) GetBookingByID(ctx context.Context, id uuid.UUID) (*entity.Booking, error) {
	rows := repo.db.QueryRow(
		ctx,
		`SELECT id, slot_id, user_id, status, conference_link, created_at FROM bookings WHERE id = $1`,
		id,
	)

	return scanBooking(rows)
}

// GetBookingByUserID - получить брони по id пользователя
func (repo *bookingRepository) GetBookingByUserID(ctx context.Context, userID uuid.UUID, from time.Time) ([]entity.Booking, error) {
	rows, err := repo.db.Query(
		ctx,
		`SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at
		 FROM bookings b
		 JOIN slots s ON s.id = b.slot_id
		 WHERE b.user_id = $1 AND s.start_time >= $2
		 ORDER BY s.start_time`,
		userID, from,
	)
	if err != nil {
		slog.Error("Error getting bookings by user",
			slog.Any("userID", userID),
			slog.Any("from", from),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("error getting bookings by user: %w", err)
	}
	defer rows.Close()

	return scanBookings(rows)
}

// ListAll - получает список броней с пагинацией
func (repo *bookingRepository) ListAll(ctx context.Context, offset, limit int) ([]entity.Booking, int, error) {
	var total int
	err := repo.db.QueryRow(ctx, `SELECT COUNT(*) FROM bookings`).Scan(&total)
	if err != nil {
		slog.Error("Error count bookings",
			slog.Any("offset", offset),
			slog.Any("limit", limit),
			slog.Any("err", err),
		)
		return nil, 0, fmt.Errorf("count bookings: %w", err)
	}

	rows, err := repo.db.Query(ctx,
		`SELECT id, slot_id, user_id, status, conference_link, created_at
		 FROM bookings ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		slog.Error("Error getting list bookings",
			slog.Any("offset", offset),
			slog.Any("limit", limit),
			slog.Any("error", err),
		)
		return nil, 0, fmt.Errorf("list bookings: %w", err)
	}
	defer rows.Close()

	bookings, err := scanBookings(rows)
	if err != nil {
		slog.Error("Error scan bookings",
			slog.Any("offset", offset),
			slog.Any("limit", limit),
			slog.Any("error", err),
		)
		return nil, 0, err
	}

	return bookings, total, nil
}

func (repo *bookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.BookingStatus) (*entity.Booking, error) {
	row := repo.db.QueryRow(
		ctx,
		`UPDATE bookings SET status = $1 WHERE id = $2
		 RETURNING id, slot_id, user_id, status, conference_link, created_at`,
		status, id,
	)
	b, err := scanBooking(row)
	if err != nil {
		slog.Error("Error scan updating booking",
			slog.Any("id", id),
			slog.Any("status", status),
			slog.Any("err", err),
		)
		return nil, err
	}

	return b, nil
}
