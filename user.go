package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type User struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	ProfilePhoto string `json:"profilePhoto"`
}

const CookieSessionName = "SCDSESSIONID"

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionHandler struct {
	sessions  map[string]uint
	userStore *UserStore
}

type UserStore struct {
	users  map[uint]*User
	mu     sync.RWMutex
	nextID uint
}

func NewUserStore() *UserStore {
	return &UserStore{
		mu: sync.RWMutex{},
		users: map[uint]*User{
			1: {
				"test@testmail.ru",
				"testpassword",
				"testuser",
				"testuser",
				"defaultphoto",
			},
		},
		nextID: 2,
	}
}

func (user User) IsValid() error {
	if user.Email == "" {
		err := errors.New("Email field is empty")
		return err
	}
	ok, _ := regexp.MatchString(`\w+[@]\w+[.]\w+`, user.Email)
	if !ok {
		err := errors.New("Email field invalid")
		return err
	}
	if user.Password == "" {
		err := errors.New("Password field is empty")
		return err
	}
	if user.FirstName == "" {
		err := errors.New("FirstName field is empty")
		return err
	}
	if user.LastName == "" {
		err := errors.New("LastName field is empty")
		return err
	}

	return nil
}

func (uStore *UserStore) AddUser(in *User) (uint, error) {
	err := in.IsValid()
	if err != nil {
		return 0, err
	}

	uStore.mu.Lock()
	id := uStore.nextID
	uStore.users[id] = in
	uStore.nextID++
	uStore.mu.Unlock()

	return id, nil
}

func (uStore *UserStore) GetUserByID(id uint) *User {
	return uStore.users[id]
}

func (uStore *UserStore) GetUserByEmail(email string) (uint, *User, bool) {
	for key, value := range uStore.users {
		if value.Email == email {
			return key, value, true
		}
	}
	return 0, nil, false
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
			Name:    CookieSessionName,
			Value:   SID,
			Expires: time.Now().Add(10 * time.Hour),
		}
		http.SetCookie(w, cookie)
	}

	if r.Method == http.MethodDelete {
		session, err := r.Cookie(CookieSessionName)
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

func (api *SessionHandler) UserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		decoder := json.NewDecoder(r.Body)
		userInput := new(User)
		err := decoder.Decode(userInput)
		if err != nil {
			log.Printf("Error while unmarshalling JSON: %s", err)
			http.Error(w, `Invalid JSONData`, 400)
			return
		}

		fmt.Println(userInput)
		var id uint
		id, err = api.userStore.AddUser(userInput)
		if err != nil {
			log.Printf("Error while Add user: %s", err)
			http.Error(w, err.Error(), 400)
			return
		}

		SID := GenerateSessionCookie()

		api.sessions[SID] = id

		cookie := &http.Cookie{
			Name:    CookieSessionName,
			Value:   SID,
			Expires: time.Now().Add(10 * time.Hour),
		}
		http.SetCookie(w, cookie)
		for key, value := range api.userStore.users {
			fmt.Println(key, value)
		}
	}
}
