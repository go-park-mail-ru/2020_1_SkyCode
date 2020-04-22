package models

import "time"

type Review struct {
	ID           uint64    `json:"id"`
	RestID       uint64    `json:"rest_id"`
	Text         string    `json:"text"`
	Author       *User     `json:"author"`
	CreationDate time.Time `json:"date"`
	Rate         float64   `json:"rate"`
}
