package models

type OrderProduct struct {
	ID        uint64 `json:"id"`
	OrderID   uint64 `json:"orderId"`
	ProductID uint64 `json:"productId" binding:"required"`
	Count     uint16 `json:"count" binding:"required"`
}

type Order struct {
	ID        uint64          `json:"id"`
	UserID    uint64          `json:"user_id"`
	Address   string          `json:"address"`
	Phone     string          `json:"phone"`
	Comment   string          `json:"comment"`
	PersonNum uint16          `json:"person_num"`
	Products  []*OrderProduct `json:"products"`
	Price     float32         `json:"price"`
}
