package models

import (
	"errors"
	"regexp"
	"sync"
)

type User struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	ProfilePhoto string `json:"profilePhoto"`
}

type UserStore struct {
	Users  map[uint]*User
	mu     sync.RWMutex
	nextID uint
}

func NewUserStore() *UserStore {
	return &UserStore{
		mu: sync.RWMutex{},
		Users: map[uint]*User{
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
	ok, _ := regexp.MatchString(`^\w+[@]\w+[.]\w+$`, user.Email)
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

func (uStore *UserStore) GetUserByID(id uint) *User {
	return uStore.Users[id]
}

func (uStore *UserStore) GetUserByEmail(email string) (uint, *User, bool) {
	for key, value := range uStore.Users {
		if value.Email == email {
			return key, value, true
		}
	}
	return 0, nil, false
}

func (uStore *UserStore) AddUser(in *User) (uint, error) {
	err := in.IsValid()
	if err != nil {
		return 0, err
	}

	uStore.mu.Lock()
	id := uStore.nextID
	uStore.Users[id] = in
	uStore.nextID++
	uStore.mu.Unlock()

	return id, nil
}
