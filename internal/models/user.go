package models

import (
	"sync"
)

type User struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Avatar    string `json:"profile_photo"`
}

type UserStore struct {
	Users  map[uint]*User
	mu     sync.RWMutex
	nextID uint
}
