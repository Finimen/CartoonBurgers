package services

import (
	"CartoonBurgers/models"
	"CartoonBurgers/repositories"
	"context"
)

type MenuService struct {
	repo *repositories.ProductRerository
}

func NewMenuService(repo *repositories.ProductRerository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) GetMenu(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return products, err
}
