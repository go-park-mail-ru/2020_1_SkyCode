package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const cookieSessionName = "SCDSESSIONID"

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionHandler struct {
	sessions  map[string]uint
	userStore *UserStore
}

type RestaurantsHandler struct {
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func GenerateSessionCookie() string {
	byteSlice := make([]rune, 64)
	for i := range byteSlice {
		byteSlice[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(byteSlice)
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

		id, user, ok := api.userStore.GetUserByEmail(loginInput.Email)
		if !ok {
			http.Error(w, `User with this email doesn't exist`, 400)
			return
		}

		if user.Password != loginInput.Password {
			http.Error(w, `Invalid password`, 400)
			return
		}

		SID := GenerateSessionCookie()

		api.sessions[SID] = id

		cookie := &http.Cookie{
			Name:    cookieSessionName,
			Value:   SID,
			Expires: time.Now().Add(10 * time.Hour),
		}
		http.SetCookie(w, cookie)
	}

	if r.Method == http.MethodDelete {
		session, err := r.Cookie(cookieSessionName)
		if err == http.ErrNoCookie {
			http.Error(w, `Access denied`, 401)
			return
		}

		_, ok := api.sessions[session.Value]
		if !ok {
			http.Error(w, `Access denied`, 401)
			return
		}

		delete(api.sessions, session.Value)

		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(w, session)
	}
}

func DefaultHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Напиши для моего пути обработчик"))
}

func main() {
	router := mux.NewRouter()

	apiSession := &SessionHandler{
		userStore: NewUserStore(),
		sessions:  make(map[string]uint, 10),
	}

	router.HandleFunc("/session", apiSession.SessionHandle)
	router.HandleFunc("/user", DefaultHandle)
	router.HandleFunc("/profile", DefaultHandle)
	router.HandleFunc("/restaurants", DefaultHandle)
	router.HandleFunc("/restaurants/{restaurant_id:[0-9]+}", DefaultHandle)

	fmt.Println("Server started")

	http.ListenAndServe(":8080", router)
}
