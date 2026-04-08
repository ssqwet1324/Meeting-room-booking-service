//go:build integration

package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestE2E_RoomScheduleBooking: переговорка → расписание → регистрация пользователя → слоты → бронь.
func TestE2E_RoomScheduleBooking(t *testing.T) {
	client := newHTTPClient()
	base := baseURL()
	waitHealth(t, client, base)

	adminTok := dummyLogin(t, client, base, "admin")
	roomName := fmt.Sprintf("e2e-room-%d", time.Now().UnixNano())
	roomID := createRoom(t, client, base, adminTok, roomName)
	createSchedule(t, client, base, adminTok, roomID)

	email := fmt.Sprintf("e2e-%d@example.com", time.Now().UnixNano())
	pass := "test-password-9"
	registerUser(t, client, base, email, pass, "user")
	userTok := login(t, client, base, email, pass)

	date := pickSlotDateUTC(t)
	ids := listSlots(t, client, base, adminTok, roomID, date)
	require.NotEmpty(t, ids, "должен быть хотя бы один слот на дату %s", date)

	_ = createBooking(t, client, base, userTok, ids[0])
}

// TestE2E_CancelBooking: полный путь + отмена брони пользователем.
func TestE2E_CancelBooking(t *testing.T) {
	client := newHTTPClient()
	base := baseURL()
	waitHealth(t, client, base)

	adminTok := dummyLogin(t, client, base, "admin")
	roomID := createRoom(t, client, base, adminTok, fmt.Sprintf("e2e-cancel-%d", time.Now().UnixNano()))
	createSchedule(t, client, base, adminTok, roomID)

	email := fmt.Sprintf("e2e-cancel-%d@example.com", time.Now().UnixNano())
	pass := "test-password-9"
	registerUser(t, client, base, email, pass, "user")
	userTok := login(t, client, base, email, pass)

	date := pickSlotDateUTC(t)
	ids := listSlots(t, client, base, adminTok, roomID, date)
	require.NotEmpty(t, ids)

	bid := createBooking(t, client, base, userTok, ids[0])
	cancelBooking(t, client, base, userTok, bid)
}
