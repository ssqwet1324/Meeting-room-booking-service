package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDummyLogin_JSON_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := newTestHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"role":"admin"}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/dummyLogin", body)
	c.Request.Header.Set("Content-Type", "application/json")

	h.DummyLogin(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["token"])
}

func TestDummyLogin_InvalidRoleInBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := newTestHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(`{"role":"god"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.DummyLogin(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_Validation_EmptyEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := newTestHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(
		`{"email":"","password":"x","role":"user"}`,
	))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Validation_MissingPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h, _ := newTestHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(
		`{"email":"a@b.com","password":""}`,
	))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
