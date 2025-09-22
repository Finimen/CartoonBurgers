package ports

import (
	"CartoonBurgers/models"
	"context"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Save(ctx context.Context, user *models.User) error
	Exists(ctx context.Context, username string) (bool, error)
}

type ProductRepository interface {
	FindAll(ctx context.Context) ([]models.Product, error)
	FindByID(ctx context.Context, id int) (*models.Product, error)
}
