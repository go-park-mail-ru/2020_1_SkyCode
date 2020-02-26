package handlers

import (
	"encoding/json"
	_models "github.com/2020_1_Skycode/internal/models"
	"net/http"
	"time"
)

func (api *SessionHandler) UserHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	userInput := new(_models.User)
	err := decoder.Decode(userInput)

	if err != nil {
		HttpResponseBody(w, "Invalid JSONData", 400)
		return
	}

	id, err := api.UserStore.AddUser(userInput)
	if err != nil {
		HttpResponseBody(w, err.Error(), 400)
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

	profile := Profile{
		Email:        userInput.Email,
		FirstName:    userInput.FirstName,
		LastName:     userInput.LastName,
		ProfilePhoto: userInput.ProfilePhoto,
	}


	_ = json.NewEncoder(w).Encode(profile)
}
