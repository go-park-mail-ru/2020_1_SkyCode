package models

type Product struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Image string  `json:"image"`
}
