package handlers

import (
	"bytes"
	json2 "encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	_handlers "github.com/2020_1_Skycode/internal/handlers"
	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/stretchr/testify/require"
)

func GetUserProfile(t *testing.T) {
	t.Parallel()

	h := _handlers.SessionHandler{
		Sessions:  make(map[string]uint, 10),
		UserStore: _models.NewUserStore(),
	}

	user := h.UserStore.GetUserByID(1)
	expected, err := json2.Marshal(user)

	if err != nil {
		t.Error("Error marshaling user")
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
	require.Equal(t, h.Sessions[cookieSession], uint(1))

	r2, _ := http.NewRequest("GET", "/profile", nil)

	r2.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	w = httptest.NewRecorder()

	h.GetUserProfile(w, r2)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
	}

	result, errRead := ioutil.ReadAll(w.Result().Body)
	if errRead != nil {
		t.Error("Error with read response body")
	}

	require.EqualValues(t, expected, result)
}
