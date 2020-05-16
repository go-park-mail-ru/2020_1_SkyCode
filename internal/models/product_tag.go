package models

type ProductTag struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	RestID uint64 `json:"rest_id"`
}
