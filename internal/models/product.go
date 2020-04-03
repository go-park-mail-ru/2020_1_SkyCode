package models

type Product struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	Price  float32 `json:"price"`
	Image  string  `json:"image"`
	RestId uint64  `json:"rest_id"`
}
