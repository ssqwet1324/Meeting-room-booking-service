package handler_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"avito/internal/entity"
	"avito/internal/handler"
	"avito/internal/middleware"
	"avito/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newBookingHandler(t *testing.T, slotID uuid.UUID, slotStart time.Time) *handler.Handler {
	t.Helper()
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "secret")
	slotRepo := &stubSlotRepoBooking{
		slot: &entity.Slot{ID: slotID, Start: slotStart},
	}
	bookings := usecase.NewBookingService(&stubBookingRepo{}, slotRepo, usecase.NewConferenceUseCase())
	return handler.New(auth, &usecase.RoomUseCase{}, &usecase.ScheduleUseCase{}, &usecase.SlotUseCase{}, bookings)
}

func TestCreateBooking_OK_User(t *testing.T) {
	gin.SetMode(gin.TestMode)
	slotID := uuid.New()
	h := newBookingHandler(t, slotID, time.Now().Add(time.Hour))
	userID := uuid.New()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleUser)
	c.Set(middleware.ContextKeyUserID, userID)
	body := fmt.Sprintf(`{"slotId":"%s","createConferenceLink":false}`, slotID)
	c.Request = httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateBooking(c)
	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "booking")
}

func TestCreateBooking_Forbidden_Admin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newBookingHandler(t, uuid.New(), time.Now().Add(time.Hour))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleAdmin)
	c.Set(middleware.ContextKeyUserID, uuid.New())
	c.Request = httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewBufferString(`{"slotId":"`+uuid.New().String()+`"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateBooking(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestListBookings_OK_Admin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "secret")
	bid := uuid.New()
	bookingsUC := usecase.NewBookingService(
		&stubBookingRepo{all: []entity.Booking{{ID: bid}}, total: 1},
		&stubSlotRepoBooking{},
		usecase.NewConferenceUseCase(),
	)
	h := handler.New(auth, &usecase.RoomUseCase{}, &usecase.ScheduleUseCase{}, &usecase.SlotUseCase{}, bookingsUC)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleAdmin)
	c.Request = httptest.NewRequest(http.MethodGet, "/bookings/list?page=1&pageSize=10", nil)

	h.ListBookings(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pagination")
}

func TestMyBookings_OK_User(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "secret")
	uid := uuid.New()
	bookingsUC := usecase.NewBookingService(
		&stubBookingRepo{all: []entity.Booking{{UserID: uid}}},
		&stubSlotRepoBooking{},
		usecase.NewConferenceUseCase(),
	)
	h := handler.New(auth, &usecase.RoomUseCase{}, &usecase.ScheduleUseCase{}, &usecase.SlotUseCase{}, bookingsUC)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleUser)
	c.Set(middleware.ContextKeyUserID, uid)
	c.Request = httptest.NewRequest(http.MethodGet, "/bookings/my", nil)

	h.MyBookings(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCancelBooking_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := usecase.NewAuthUseCase(&stubUserRepo{}, "secret")
	bid := uuid.New()
	uid := uuid.New()
	bookingsUC := usecase.NewBookingService(
		&stubBookingRepo{booking: &entity.Booking{ID: bid, UserID: uid, Status: entity.BookingStatusActive}},
		&stubSlotRepoBooking{},
		usecase.NewConferenceUseCase(),
	)
	h := handler.New(auth, &usecase.RoomUseCase{}, &usecase.ScheduleUseCase{}, &usecase.SlotUseCase{}, bookingsUC)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleUser)
	c.Set(middleware.ContextKeyUserID, uid)
	c.Params = gin.Params{{Key: "bookingId", Value: bid.String()}}
	c.Request = httptest.NewRequest(http.MethodPost, "/cancel", nil)

	h.CancelBooking(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
