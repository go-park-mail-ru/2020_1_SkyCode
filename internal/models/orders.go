package models

type OrderProduct struct {
	ID        uint64 `json:"id,omitempty"`
	OrderID   uint64 `json:"orderId,omitempty"`
	ProductID uint64 `json:"productId" binding:"required"`
	Count     uint32 `json:"count" binding:"required"`
}

type Order struct {
	ID        uint64     `json:"id"`
	UserID    uint64     `json:"user_id"`
	RestID    uint64     `json:"restId"`
	RestName  string     `json:"rest_name"`
	Address   string     `json:"address"`
	Phone     string     `json:"phone"`
	Comment   string     `json:"comment"`
	PersonNum uint32     `json:"person_num"`
	Products  []*Product `json:"products"`
	Price     float32    `json:"price"`
	CreatedAt string     `json:"created_at"`
	Status    string     `json:"status"`
}
