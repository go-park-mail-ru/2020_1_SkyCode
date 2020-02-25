package handlers

import (
	"bytes"
	json2 "encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_handlers "github.com/2020_1_Skycode/internal/handlers"
	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/stretchr/testify/require"
)

func TestGetUserProfile(t *testing.T) {
	t.Parallel()

	h := _handlers.SessionHandler{
		Sessions:  make(map[string]uint, 10),
		UserStore: _models.NewUserStore(),
	}

	user := h.UserStore.GetUserByID(1)

	profile := _handlers.Profile{
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ProfilePhoto: user.ProfilePhoto,
	}
	expected, err := json2.Marshal(profile)

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

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	body := bytes.NewReader([]byte(
		`{
			"email": "test@testmail.ru",
			"password": "testpassword"
			}`,
	))

	r := httptest.NewRequest("POST", "/session", body)
	w := httptest.NewRecorder()

	h := _handlers.SessionHandler{
		Sessions:  make(map[string]uint, 10),
		UserStore: _models.NewUserStore(),
	}

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

	expectedProfile := _handlers.Profile{
		Email:     "newemail@mail.ru",
		FirstName: "New",
		LastName:  "Name",
	}

	bodyUpdate := &bytes.Buffer{}

	writer := multipart.NewWriter(bodyUpdate)
	writer.WriteField("email", "newemail@mail.ru")
	writer.WriteField("firstname", "New")
	writer.WriteField("lastname", "Name")

	part, _ := writer.CreateFormFile("profilephoto", "testfile")

	part.Write([]byte("SOME FILE CONTENT"))

	writer.Close()

	ru, errReq := http.NewRequest("PUT", "/profile", bodyUpdate)
	if errReq != nil {
		t.Error("Error create form request")
	}
	ru.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	ru.Header.Add("Content-Type", writer.FormDataContentType())
	wu := httptest.NewRecorder()

	h.GetUserProfile(wu, ru)

	if w.Code != http.StatusOK {
		t.Error("Status is not ok")
	}

	result, errRead := ioutil.ReadAll(wu.Result().Body)
	if errRead != nil {
		t.Error("Error with read response body")
	}

	resultProfile := _handlers.Profile{}
	json2.Unmarshal(result, &resultProfile)
	require.NotEqual(t, "", resultProfile.ProfilePhoto)
	os.RemoveAll("images/")
	resultProfile.ProfilePhoto = ""

	require.EqualValues(t, expectedProfile, resultProfile)
}
