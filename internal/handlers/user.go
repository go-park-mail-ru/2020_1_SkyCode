package handlers

import (
	"encoding/json"
	"fmt"
	_models "github.com/2020_1_Skycode/internal/models"
	"log"
	"net/http"
	"time"
)

func DefaultHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Напиши для моего пути обработчик"))
}

func (api *SessionHandler) UserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		userInput := new(_models.User)
		err := decoder.Decode(userInput)

		if err != nil {
			log.Printf("Error while unmarshalling JSON: %s", err)
			http.Error(w, `Invalid JSONData`, 400)
			return
		}

		fmt.Println(userInput)
		var id uint
		id, err = api.UserStore.AddUser(userInput)
		if err != nil {
			log.Printf("Error while Add user: %s", err)
			http.Error(w, err.Error(), 400)
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
		for key, value := range api.UserStore.Users {
			fmt.Println(key, value)
		}
	}
}
