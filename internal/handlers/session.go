package handlers

import (
	"encoding/json"
	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type SessionHandler struct {
	Sessions  map[string]uint
	UserStore *_models.UserStore
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *SessionHandler) SessionHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		loginInput := new(LoginInput)
		err := decoder.Decode(loginInput)

		if err != nil {
			HttpResponseBody(w, "Bad params", 400)
			return
		}

		id, user, ok := api.UserStore.GetUserByEmail(loginInput.Email)
		if !ok {
			HttpResponseBody(w, "No such user", 400)
			return
		}

		if user.Password != loginInput.Password {
			HttpResponseBody(w, "Invalid password", 400)
			return
		}

		SID := uuid.New().String()
		api.Sessions[SID] = id

		cookie := &http.Cookie{
			Name:    _models.CookieSessionName,
			Value:   SID,
			Expires: time.Now().Add(10 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		_ = json.NewEncoder(w).Encode(user)

		return
	}

	if r.Method == http.MethodDelete {
		session, err := r.Cookie(_models.CookieSessionName)
		if err == http.ErrNoCookie {
			HttpResponseBody(w, "Access denied", 401)
			return
		}

		_, ok := api.Sessions[session.Value]
		if !ok {
			HttpResponseBody(w, "Access denied", 401)
			return
		}

		delete(api.Sessions, session.Value)

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)

		HttpResponseBody(w, "", 200)
	}
}
