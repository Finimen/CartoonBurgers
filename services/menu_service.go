package services

import (
	"CartoonBurgers/models"
	"CartoonBurgers/repositories"
	"context"
)

type MenuService struct {
	repo *repositories.MenuPerository
}

func NewMenuService(repo repositories.MenuPerository) *MenuService {
	return &MenuService{repo: &repo}
}

func (s *MenuService) GetMenu(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return products, err
}
