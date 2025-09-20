package models

type ProductCategory int

const (
	Burger ProductCategory = iota
	Snack
	Drink
	Dessert
)

type ProductType int

const (
	Default ProductType = iota
	Latest
)

type Product struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Price    int             `json:"price"`
	Count    int             `json:"count"`
	Type     ProductType     `json:"type"`
	Category ProductCategory `json:"category"`
}
