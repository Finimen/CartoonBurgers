package models

type CartItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}
