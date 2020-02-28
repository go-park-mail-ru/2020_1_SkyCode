package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	_models "github.com/2020_1_Skycode/internal/models"
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

	_, _, userAppear := api.UserStore.GetUserByEmail(userInput.Email)
	if userAppear {
		HttpResponseBody(w, "User with this data already appears", 400)
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
		ProfilePhoto: "default.jpg",
	}

	_ = json.NewEncoder(w).Encode(profile)
}
