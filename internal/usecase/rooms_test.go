package usecase

import (
	"context"
	"testing"

	"avito/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomUseCase_ListRooms(t *testing.T) {
	rid := uuid.New()
	uc := NewRoomUseCase(&mockRoomRepo{
		rooms: []entity.Room{{ID: rid, Name: "Hall"}},
	})
	list, err := uc.ListRooms(context.Background())
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Hall", list[0].Name)
}

func TestRoomUseCase_GetRoomByID(t *testing.T) {
	rid := uuid.New()
	uc := NewRoomUseCase(&mockRoomRepo{room: &entity.Room{ID: rid, Name: "A"}})
	r, err := uc.GetRoomByID(context.Background(), rid)
	require.NoError(t, err)
	assert.Equal(t, rid, r.ID)
}

func TestRoomUseCase_CreateRoom(t *testing.T) {
	uc := NewRoomUseCase(&mockRoomRepo{})
	r, err := uc.CreateRoom(context.Background(), entity.NewRoom{Name: "New"})
	require.NoError(t, err)
	assert.Equal(t, "New", r.Name)
}
