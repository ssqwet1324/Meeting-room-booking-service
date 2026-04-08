//go:build integration

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func baseURL() string {
	if u := os.Getenv("E2E_BASE_URL"); u != "" {
		return u
	}
	return "http://127.0.0.1:8080"
}

func newHTTPClient() *http.Client {
	return &http.Client{Timeout: 60 * time.Second}
}

var (
	serverOnce sync.Once
	serverUp   bool
	serverBase string
)

// waitHealth проверяет один раз на процесс, что API доступен
func waitHealth(t *testing.T, client *http.Client, base string) {
	t.Helper()
	serverOnce.Do(func() {
		serverBase = base
		deadline := time.Now().Add(20 * time.Second)
		for time.Now().Before(deadline) {
			resp, err := client.Get(base + "/_info")
			if err == nil && resp.StatusCode == http.StatusOK {
				_ = resp.Body.Close()
				serverUp = true
				return
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
			time.Sleep(300 * time.Millisecond)
		}
		serverUp = false
	})
	if !serverUp {
		t.Skipf("E2E: сервер недоступен по %s (запустите: docker compose up -d)", serverBase)
	}
}

func doRequest(t *testing.T, client *http.Client, method, urlStr, token string, body []byte) (int, []byte) {
	t.Helper()
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, urlStr, r)
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp.StatusCode, b
}

func requireStatus(t *testing.T, want int, got int, body []byte) {
	t.Helper()
	if got != want {
		t.Fatalf("ожидался HTTP %d, получен %d, тело: %s", want, got, string(body))
	}
}

func dummyLogin(t *testing.T, client *http.Client, base, role string) string {
	t.Helper()
	payload := fmt.Sprintf(`{"role":%q}`, role)
	code, body := doRequest(t, client, http.MethodPost, base+"/dummyLogin", "", []byte(payload))
	requireStatus(t, http.StatusOK, code, body)
	var out struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.Unmarshal(body, &out))
	require.NotEmpty(t, out.Token)
	return out.Token
}

func registerUser(t *testing.T, client *http.Client, base, email, password, role string) {
	t.Helper()
	payload := fmt.Sprintf(`{"email":%q,"password":%q,"role":%q}`, email, password, role)
	code, body := doRequest(t, client, http.MethodPost, base+"/register", "", []byte(payload))
	requireStatus(t, http.StatusCreated, code, body)
}

func login(t *testing.T, client *http.Client, base, email, password string) string {
	t.Helper()
	payload := fmt.Sprintf(`{"email":%q,"password":%q}`, email, password)
	code, body := doRequest(t, client, http.MethodPost, base+"/login", "", []byte(payload))
	requireStatus(t, http.StatusOK, code, body)
	var out struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.Unmarshal(body, &out))
	require.NotEmpty(t, out.Token)
	return out.Token
}

func createRoom(t *testing.T, client *http.Client, base, adminToken, name string) string {
	t.Helper()
	payload := fmt.Sprintf(`{"name":%q}`, name)
	code, body := doRequest(t, client, http.MethodPost, base+"/rooms/create", adminToken, []byte(payload))
	requireStatus(t, http.StatusCreated, code, body)
	var out struct {
		Room struct {
			ID string `json:"id"`
		} `json:"room"`
	}
	require.NoError(t, json.Unmarshal(body, &out))
	require.NotEmpty(t, out.Room.ID)
	return out.Room.ID
}

func createSchedule(t *testing.T, client *http.Client, base, adminToken, roomID string) {
	t.Helper()
	payload := `{"daysOfWeek":[1,2,3,4,5],"startTime":"09:00","endTime":"18:00"}`
	url := fmt.Sprintf("%s/rooms/%s/schedule/create", base, roomID)
	code, body := doRequest(t, client, http.MethodPost, url, adminToken, []byte(payload))
	requireStatus(t, http.StatusCreated, code, body)
}

// pickSlotDateUTC возвращает дату YYYY-MM-DD (UTC), в который есть рабочий день и первый слот ещё не в прошлом.
func pickSlotDateUTC(t *testing.T) string {
	t.Helper()
	for day := 1; day <= 40; day++ {
		d := time.Now().UTC().AddDate(0, 0, day)
		wd := d.Weekday()
		if wd == time.Saturday || wd == time.Sunday {
			continue
		}
		dayStart := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		firstSlot := dayStart.Add(9 * time.Hour)
		if firstSlot.After(time.Now()) {
			return dayStart.Format("2006-01-02")
		}
	}
	t.Fatal("не удалось подобрать дату со слотами в будущем")
	return ""
}

func listSlots(t *testing.T, client *http.Client, base, adminToken, roomID, date string) []string {
	t.Helper()
	url := fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", base, roomID, date)
	code, body := doRequest(t, client, http.MethodGet, url, adminToken, nil)
	requireStatus(t, http.StatusOK, code, body)
	var out struct {
		Slots []struct {
			ID string `json:"id"`
		} `json:"slots"`
	}
	require.NoError(t, json.Unmarshal(body, &out))
	ids := make([]string, 0, len(out.Slots))
	for _, s := range out.Slots {
		ids = append(ids, s.ID)
	}
	return ids
}

func createBooking(t *testing.T, client *http.Client, base, userToken, slotID string) string {
	t.Helper()
	payload := fmt.Sprintf(`{"slotId":%q,"createConferenceLink":false}`, slotID)
	code, body := doRequest(t, client, http.MethodPost, base+"/bookings/create", userToken, []byte(payload))
	requireStatus(t, http.StatusCreated, code, body)
	var out struct {
		Booking struct {
			ID string `json:"id"`
		} `json:"booking"`
	}
	require.NoError(t, json.Unmarshal(body, &out))
	require.NotEmpty(t, out.Booking.ID)
	return out.Booking.ID
}

func cancelBooking(t *testing.T, client *http.Client, base, userToken, bookingID string) {
	t.Helper()
	url := fmt.Sprintf("%s/bookings/%s/cancel", base, bookingID)
	code, body := doRequest(t, client, http.MethodPost, url, userToken, nil)
	requireStatus(t, http.StatusOK, code, body)
	var out struct {
		Booking struct {
			Status string `json:"status"`
		} `json:"booking"`
	}
	require.NoError(t, json.Unmarshal(body, &out))
	require.Equal(t, "cancelled", out.Booking.Status)
}
