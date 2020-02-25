package handlers

import (
	json2 "encoding/json"
	"fmt"
	_models "github.com/2020_1_Skycode/internal/models"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
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
	if r.Method == http.MethodGet {
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
			Email:        user.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
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

	if r.Method == http.MethodPut {
		r.ParseMultipartForm(5 * 1024 * 1025)
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
		tempUser, err := api.updateUser(r, *user)

		if err != nil {
			fmt.Println(err)
			http.Error(w, `Invalid JSONData`, 400)
			return
		}

		*user = tempUser

		file, _, err := r.FormFile("profilephoto")
		defer file.Close()

		if user.ProfilePhoto != "" {
			err := os.Remove(user.ProfilePhoto)

			if err != nil {
				fmt.Println(err)
				http.Error(w, `Server error`, 500)
				return
			}
		}

		id := uuid.New()
		data, _ := ioutil.ReadAll(file)
		filePath := `images/` + id.String() + `.jpg`
		err = ioutil.WriteFile(filePath, data, 0644)

		if err != nil {
			fmt.Println(err)
			http.Error(w, `Server error`, 500)
			return
		}

		user.ProfilePhoto = filePath

		profile := Profile{
			Email:        user.Email,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			ProfilePhoto: user.ProfilePhoto,
		}

		jsonProfile, _ := json2.Marshal(profile)

		w.Write(jsonProfile)
	}
}
