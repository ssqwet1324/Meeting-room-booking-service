package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito/internal/entity"
	"avito/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListRooms_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roomID := uuid.New()
	h := newHandlerWithRooms(t, roomID)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/rooms/list", nil)

	h.ListRooms(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Hall")
}

func TestCreateRoom_Forbidden_NotAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newHandlerWithRooms(t, uuid.New())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleUser)
	c.Request = httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewBufferString(`{"name":"X"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateRoom(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreateRoom_OK_Admin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newHandlerWithRooms(t, uuid.New())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleAdmin)
	c.Request = httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewBufferString(`{"name":"Board"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateRoom(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Board")
}

func TestListSlots_MissingDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roomID := uuid.New()
	h := newHandlerWithRooms(t, roomID)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "roomId", Value: roomID.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/slots/list", nil)

	h.ListSlots(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListSlots_InvalidRoomID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newHandlerWithRooms(t, uuid.New())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "roomId", Value: "not-uuid"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/x/slots/list?date=2024-06-10", nil)

	h.ListSlots(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListSlots_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roomID := uuid.New()
	h := newHandlerWithRooms(t, roomID)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "roomId", Value: roomID.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/slots/list?date=2024-06-10", nil)

	h.ListSlots(c)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "slots")
}

func TestCreateSchedule_Forbidden_NotAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roomID := uuid.New()
	h := newHandlerWithRooms(t, roomID)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleUser)
	c.Params = gin.Params{{Key: "roomId", Value: roomID.String()}}
	c.Request = httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateSchedule(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreateSchedule_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roomID := uuid.New()
	h := newHandlerWithRooms(t, roomID)

	body := `{"daysOfWeek":[1,2,3],"startTime":"09:00","endTime":"18:00"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.ContextKeyRole, entity.RoleAdmin)
	c.Params = gin.Params{{Key: "roomId", Value: roomID.String()}}
	c.Request = httptest.NewRequest(http.MethodPost, "/create", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateSchedule(c)
	assert.Equal(t, http.StatusCreated, w.Code)
}
