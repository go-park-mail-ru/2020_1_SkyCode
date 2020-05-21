package models

import "time"

type Notification struct {
	ID           uint64
	UserID       uint64
	OrderID      uint64
	UnreadStatus bool
	Status       string
	DateTime     time.Time
}
