package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUserAndAuthorized(t *testing.T) {
	t.Parallel()

	h := SessionHandler{
		sessions:  make(map[string]uint, 10),
		userStore: NewUserStore(),
	}

	body := bytes.NewReader([]byte(
		`{
			"email": "testemail@test.com",
			"password": "password",
			"firstName": "Test",
			"lastName": "User"
			}`,
	))

	expectedUsers := map[uint]*User{
		1: {
			"test@testmail.ru",
			"testpassword",
			"testuser",
			"testuser",
			"defaultphoto",
		},
		2: {
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
	}

	cookie := w.Result().Cookies()
	if len(cookie) == 0 {
		t.Error("Cookie is not seted")
	}

	require.EqualValues(t, h.userStore.users, expectedUsers)

	cookieSession := cookie[0].Value
	if h.sessions[cookieSession] != 2 {
		t.Error("Cookie set incorrect")
	}
}

func TestAuthorizingAndLogout(t *testing.T) {
	t.Parallel()

	h := SessionHandler{
		sessions:  make(map[string]uint, 10),
		userStore: NewUserStore(),
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
	}

	cookie := w.Result().Cookies()
	if len(cookie) == 0 {
		t.Error("Cookie is not seted")
	}
	cookieSession := cookie[0].Value
	require.Equal(t, h.sessions[cookieSession], uint(1))

	r2, _ := http.NewRequest("DELETE", "/session", nil)

	r2.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	w = httptest.NewRecorder()

	h.SessionHandle(w, r2)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
	}

	require.Equal(t, h.sessions, map[string]uint{})
}
