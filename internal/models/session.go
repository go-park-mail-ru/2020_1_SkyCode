package models

import (
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"time"
)

type Session struct {
	ID uint64 `json:"id"`
	UserId uint64 `json:"userId"`
	Token string `json:"token"`
}

func GenerateSession(userId uint64) (*Session, *http.Cookie) {
	cookie := &http.Cookie{
		Name:       "SkyDelivery",
		Value:      uuid.New().String(),
		MaxAge:    time.Now().Add(24 * time.Hour).Second(),
		HttpOnly:   true,
	}

	return &Session{
		UserId: userId,
		Token:  cookie.Value,
	}, cookie
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)
const CookieSessionName = "SCDSESSIONID"

func GenerateSessionCookie() string {
	byteSlice := make([]rune, 64)
	for i := range byteSlice {
		byteSlice[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(byteSlice)
}