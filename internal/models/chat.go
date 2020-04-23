package models

import "time"

type Chat struct {
	UserID uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	ChatID   string `json:"chat_id"`
}

type ChatMessage struct {
	UserID uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	ChatID string `json:"chat_id"`
	Message string `json:"message"`
	Created time.Time `json:"created"`
}


