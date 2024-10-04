package models

type Payment struct {
	OrderId int  `json:"OrderID"`
	Status  bool `json:"Status"`
}
