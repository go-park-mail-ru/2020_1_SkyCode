package models

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Session struct {
	ID uint64 `json:"id"`
	UserId uint64 `json:"userId"`
	Token string `json:"token"`
	Expiration time.Time `json:"expiration, omitempty"`
}

func GenerateSession(userId uint64) (*Session, *http.Cookie) {
	cookie := &http.Cookie{
		Name:       "SkyDelivery",
		Value:      uuid.New().String(),
		MaxAge:    3600 * 12,
		HttpOnly:   true,
		SameSite: http.SameSiteNoneMode,
	}

	return &Session{
		UserId: userId,
		Token:  cookie.Value,
	}, cookie
}