package usecase

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConferenceUseCase_CreateURL(t *testing.T) {
	uc := NewConferenceUseCase()
	id := uuid.MustParse("00000000-0000-0000-0000-000000000099")
	link, err := uc.CreateURL(id)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(link, "https://example.com/meeting/"))
	assert.Contains(t, link, id.String())
}
