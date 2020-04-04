package models

type OrderProduct struct {
	ID        uint64 `json:"id"`
	OrderID   uint64 `json:"order_id"`
	ProductID uint64 `json:"product_id"`
	Count     uint16 `json:"count"`
}

type Order struct {
	ID        uint64          `json:"id"`
	UserID    uint64          `json:"user_id"`
	Address   string          `json:"address"`
	Comment   string          `json:"comment"`
	PersonNum uint16          `json:"person_num"`
	Products  []*OrderProduct `json:"products"`
	Price     float32         `json:"price"`
}
