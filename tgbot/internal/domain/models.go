package domain

import "time"

type Status string
type ShopState string

const (
	NONE                            ShopState = ""
	WAIT_PAY_ANSWER                 ShopState = "WAIT_PAY_ANSWER"
	WAIT_ADD_TO_SHOPPING_CART       ShopState = "WAIT_ADD_TO_SHOPPING_CART"
	WAIT_ORDER_ID_TO_GET_ORDER_INFO ShopState = "WAIT_ORDER_ID_TO_GET_ORDER_INFO"
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
	ShopCartSt  ShopCartState  `json:"shop_cart_st"`
}

type ProductState struct {
	CurPage int64 `json:"cur_page"`
}

type ShopOrderState struct {
	State    ShopState `json:"state"`
	TryCount int64     `json:"try_count"`
}

type ShopCartState struct {
	Products map[int64]OrderProduct
}

func NewUserState(userID int64) UserState {
	userState := UserState{}
	userState.UserID = userID
	userState.ShopCartSt.Products = make(map[int64]OrderProduct, 10)
	return userState
}
