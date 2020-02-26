package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	_handlers "github.com/2020_1_Skycode/internal/handlers"
	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAndAuthorized(t *testing.T) {
	t.Parallel()

	h := _handlers.SessionHandler{
		Sessions:  make(map[string]uint, 10),
		UserStore: _models.NewUserStore(),
	}

	body := bytes.NewReader([]byte(
		`{
			"email": "testemail@test.com",
			"password": "password",
			"firstName": "Test",
			"lastName": "User"
			}`,
	))

	expectedUsers := map[uint]*_models.User{
		1: {
			"test@testmail.ru",
			"testpassword",
			"testuser",
			"testuser",
			"",
		},
		2: {
			"t@m.ru",
			"pass",
			"testuser",
			"testuser",
			"",
		},
		3: {
			"testemail@test.com",
			"password",
			"Test",
			"User",
			"",
		},
	}

	r := httptest.NewRequest("POST", "/user", body)
	w := httptest.NewRecorder()

	h.UserHandle(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	cookie := w.Result().Cookies()
	if len(cookie) == 0 {
		t.Error("Cookie is not seted")
		return
	}

	require.EqualValues(t, expectedUsers, h.UserStore.Users)

	cookieSession := cookie[0].Value
	if h.Sessions[cookieSession] != 3 {
		t.Error("Cookie set incorrect")
		return
	}
}

func TestAuthorizingAndLogout(t *testing.T) {
	t.Parallel()

	h := _handlers.SessionHandler{
		Sessions:  make(map[string]uint, 10),
		UserStore: _models.NewUserStore(),
	}

	body := bytes.NewReader([]byte(
		`{
			"email": "test@testmail.ru",
			"password": "testpassword"
			}`,
	))

	r := httptest.NewRequest("POST", "/session", body)
	w := httptest.NewRecorder()

	h.SessionHandle(w, r)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	cookie := w.Result().Cookies()
	if len(cookie) == 0 {
		t.Error("Cookie is not seted")
		return
	}
	cookieSession := cookie[0].Value
	require.Equal(t, h.Sessions[cookieSession], uint(1))

	r2, _ := http.NewRequest("DELETE", "/session", nil)

	r2.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	w = httptest.NewRecorder()

	h.SessionHandle(w, r2)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
		return
	}

	require.Equal(t, h.Sessions, map[string]uint{})
}
