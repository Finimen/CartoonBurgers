package ports

import (
	"CartoonBurgers/models"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (*AuthResult, error)
	Register(ctx context.Context, user *models.User) error
	Logout(ctx context.Context, token string) error
}

type AuthResult struct {
	Token string
	User  *models.User
}

type CartService interface {
	GetCart(ctx context.Context, userID string) ([]models.CartItem, error)
	AddToCart(ctx context.Context, userID string, item models.CartItem) error
	RemoveFromCart(ctx context.Context, userID string, productID int) error
}
