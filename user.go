package main

import "sync"

type User struct {
	Email        string `json:"username"`
	Password     string `json:"password"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	ProfilePhoto string `json:"profilePhoto"`
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

func (uStore *UserStore) AddUser(in *User) (uint, error) {
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
