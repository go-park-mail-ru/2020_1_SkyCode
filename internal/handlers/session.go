package handlers

import (
	"encoding/json"
	"fmt"
	_models "github.com/2020_1_Skycode/internal/models"
	"log"
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
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		loginInput := new(LoginInput)
		err := decoder.Decode(loginInput)
		if err != nil {
			log.Printf("Error while unmarshalling JSON: %s", err)
			http.Error(w, `Invalid JSONData`, 400)
			return
		}

		id, user, ok := api.UserStore.GetUserByEmail(loginInput.Email)
		if !ok {
			http.Error(w, `User with this email doesn't exist`, 400)
			return
		}

		if user.Password != loginInput.Password {
			http.Error(w, `Invalid password`, 400)
			return
		}

		SID := _models.GenerateSessionCookie()

		api.Sessions[SID] = id

		cookie := &http.Cookie{
			Name:    _models.CookieSessionName,
			Value:   SID,
			Expires: time.Now().Add(10 * time.Hour),
		}
		http.SetCookie(w, cookie)
	}

	if r.Method == http.MethodDelete {
		session, err := r.Cookie(_models.CookieSessionName)
		if err == http.ErrNoCookie {
			http.Error(w, `Access denied`, 401)
			return
		}

		_, ok := api.Sessions[session.Value]
		if !ok {
			http.Error(w, `Access denied`, 401)
			return
		}

		delete(api.Sessions, session.Value)

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	}

	fmt.Println(api.Sessions)
}
