package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/google/uuid"
)

type Profile struct {
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	ProfilePhoto string `json:"profilePhoto"`
}

func (api *SessionHandler) isUniqueUserEmail(email string) bool {
	for _, user := range api.UserStore.Users {
		if user.Email == email {
			return false
		}
	}

	return true
}

func (api *SessionHandler) updateUser(r *http.Request, user _models.User) (_models.User, error) {
	email := r.FormValue("email")
	if email != "" {
		if api.isUniqueUserEmail(email) {
			user.Email = email
		}
	}

	firstName := r.FormValue("firstname")
	if firstName != "" {
		user.FirstName = firstName
	}

	lastName := r.FormValue("lastname")
	if lastName != "" {
		user.LastName = lastName
	}

	password := r.FormValue("password")
	if password != "" {
		user.Password = password
	}

	return user, user.IsValid()
}

func (api *SessionHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		session, err := r.Cookie(_models.CookieSessionName)

		if session == nil || err == http.ErrNoCookie {
			HttpResponseBody(w, "Unauthorized", 401)
			return
		}

		userID, ok := api.Sessions[session.Value]
		if !ok {
			HttpResponseBody(w, "Unauthorized", 401)
			return
		}

		user := api.UserStore.Users[userID]

		profile := Profile{
			Email:        user.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			ProfilePhoto: user.ProfilePhoto,
		}

		_ = json.NewEncoder(w).Encode(profile)

		return
	}

	if r.Method == http.MethodPut {
		err := r.ParseMultipartForm(5 * 1024 * 1025)

		if err != nil {
			HttpResponseBody(w, "Request Entity Too Large", 413)
			return
		}

		session, err := r.Cookie(_models.CookieSessionName)

		if session == nil {
			HttpResponseBody(w, "Server error", 500)
			return
		}

		userID, ok := api.Sessions[session.Value]
		if !ok {
			HttpResponseBody(w, "Unauthorized", 401)
			return
		}

		user := api.UserStore.Users[userID]

		tempUser, err := api.updateUser(r, *user)

		if err != nil {
			HttpResponseBody(w, "Bad request", 400)
			return
		}

		*user = tempUser

		file, _, err := r.FormFile("profilephoto")
		if file != nil {
			defer file.Close()
			if user.ProfilePhoto != "" {
				err := os.Remove("images/" + user.ProfilePhoto)

				if err != nil {
					HttpResponseBody(w, err.Error(), 500)
					return
				}
			}

			id := uuid.New()
			data, _ := ioutil.ReadAll(file)
			filePath := id.String() + `.jpg`

			if _, err := os.Stat("images/"); os.IsNotExist(err) {
				err := os.Mkdir("images", 0775)

				if err != nil {
					HttpResponseBody(w, err.Error(), 500)
					return
				}
			}

			err = ioutil.WriteFile("images/" + filePath, data, 0644)

			if err != nil {
				HttpResponseBody(w, err.Error(), 500)
				return
			}

			user.ProfilePhoto = filePath
		}

		profile := Profile{
			Email:        user.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			ProfilePhoto: user.ProfilePhoto,
		}

		_ = json.NewEncoder(w).Encode(profile)

		return
	}
}
