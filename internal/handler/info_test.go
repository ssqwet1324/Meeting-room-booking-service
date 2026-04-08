package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"avito/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Info(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := handler.New(nil, nil, nil, nil, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h.Info(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}
