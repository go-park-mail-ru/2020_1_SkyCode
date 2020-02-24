package handlers

import (
	json2 "encoding/json"
	"fmt"
	_models "github.com/2020_1_Skycode/internal/models"
	"net/http"
)

type Profile struct {
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	ProfilePhoto string `json:"profilePhoto"`
}

func (api *SessionHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `No such method`, 405)
		return
	}

	session, err := r.Cookie(_models.CookieSessionName)
	if err == http.ErrNoCookie {
		http.Error(w, `Access denied`, 401)
		return
	}

	userID, ok := api.Sessions[session.Value]
	if !ok {
		http.Error(w, `Access denied`, 401)
		return
	}

	user := api.UserStore.Users[userID]

	profile := Profile{
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		ProfilePhoto: user.ProfilePhoto,
	}

	jsonProfile, err := json2.Marshal(profile)

	if err != nil {
		fmt.Println(err)
		http.Error(w, `Server error`, 500)
		return
	}

	_, err = w.Write(jsonProfile)

	if err != nil {
		fmt.Println(err)
		http.Error(w, `Server error`, 500)
	}
}
