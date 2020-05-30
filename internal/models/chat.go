package models

import "time"

type Chat struct {
	UserID        uint64 `json:"user_id"`
	UserName      string `json:"user_name"`
	ChatID        uint64 `json:"chat_id"`
	UserConnected bool   `json:"user_connected"`
	SupConnected  bool   `json:"sup_connected"`
}

type ChatMessage struct {
	UserID   uint64    `json:"user_id"`
	UserName string    `json:"user_name"`
	ChatID   uint64    `json:"chat_id"`
	Message  string    `json:"message"`
	Created  time.Time `json:"created"`
}
