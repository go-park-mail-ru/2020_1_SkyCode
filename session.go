package main

import "sync"

type Session struct {
	ID     uint
	UserID uint
}

type SessionStore struct {
	sessions []*Session
	mu       sync.RWMutex
	nextID   uint
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		mu:       sync.RWMutex{},
		sessions: []*Session{},
	}
}
