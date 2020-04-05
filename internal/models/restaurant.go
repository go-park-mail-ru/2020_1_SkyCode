package models

type Restaurant struct {
	ID          uint64     `json:"id"`
	ManagerID   uint64     `json:"managerId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Rating      float32    `json:"rating"`
	Image       string     `json:"image"`
	Products    []*Product `json:"products"`
}

type ResStorage struct {
	Restaurants map[uint]*Restaurant
	nextID      uint
}
