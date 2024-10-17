package orders

import "time"

type Product struct {
	ArticleID int64
	Count     int64
}

type OrderInfo struct {
	OrderId     int64
	Description string
	State       string
}

type Order struct {
	OrderInfo
	OrderProducts []Product
}

type OrderStatus struct {
	Status string
	Time   time.Time
}
