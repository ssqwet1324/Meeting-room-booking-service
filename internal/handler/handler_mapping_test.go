package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"avito/internal/entity"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCodeToStatus(t *testing.T) {
	assert.Equal(t, http.StatusBadRequest, codeToStatus(entity.CodeInvalidRequest))
	assert.Equal(t, http.StatusUnauthorized, codeToStatus(entity.CodeUnauthorized))
	assert.Equal(t, http.StatusForbidden, codeToStatus(entity.CodeForbidden))
	assert.Equal(t, http.StatusNotFound, codeToStatus(entity.CodeNotFound))
	assert.Equal(t, http.StatusNotFound, codeToStatus(entity.CodeRoomNotFound))
	assert.Equal(t, http.StatusNotFound, codeToStatus(entity.CodeSlotNotFound))
	assert.Equal(t, http.StatusNotFound, codeToStatus(entity.CodeBookingNotFound))
	assert.Equal(t, http.StatusConflict, codeToStatus(entity.CodeSlotBooked))
	assert.Equal(t, http.StatusConflict, codeToStatus(entity.CodeScheduleExists))
	assert.Equal(t, http.StatusInternalServerError, codeToStatus(entity.CodeInternalError))
	assert.Equal(t, http.StatusInternalServerError, codeToStatus(entity.Code("UNKNOWN")))
}

func TestMapError_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mapError(c, entity.New(entity.CodeInvalidRequest, entity.ErrMsgInvalidRequest))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMapError_NonAppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mapError(c, assert.AnError)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestWriteJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	writeJSON(c, http.StatusOK, gin.H{"a": 1})
	assert.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"a"`)
}
