package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppError_Error(t *testing.T) {
	e := New(CodeInvalidRequest, "bad request")
	assert.Equal(t, "bad request", e.Error())
}

func TestGetCode_FromAppError(t *testing.T) {
	code, ok := GetCode(New(CodeUnauthorized, "nope"))
	require.True(t, ok)
	assert.Equal(t, CodeUnauthorized, code)
}

func TestGetCode_NotAppError(t *testing.T) {
	_, ok := GetCode(errors.New("plain"))
	assert.False(t, ok)
}
