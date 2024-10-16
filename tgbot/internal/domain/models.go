package domain

import "time"

type Status string
type ShopState string

const (
	NONE           ShopState = ""
	WAIT_ADD_ORDER ShopState = "WAIT_ADD_ORDER"
)

type Product struct {
	ArticleID   int64
	Category    string
	Name        string
	Description string
	Count       int64
}

type CreateOrderInfo struct {
	UserID      int64
	Description string
	Products    []OrderProduct
}

type OrderProduct struct {
	ArticleID int64
	Count     int64
}

type Order struct {
	OrderID       int64
	Products      []Product
	Description   string
	CurrentStatus string
	StatusHistory []StatusInfo
}

type StatusInfo struct {
	Status string
	Time   time.Time
}

type UserState struct {
	UserID      int64          `json:"user_id"`
	ProductSt   ProductState   `json:"product_state"`
	ShopOrderSt ShopOrderState `json:"shop_order_state"`
}

type ProductState struct {
	CurPage  int64 `json:"cur_page"`
	Products map[int64]OrderProduct
}

type ShopOrderState struct {
	State    ShopState `json:"state"`
	TryCount int64     `json:"try_count"`
}
