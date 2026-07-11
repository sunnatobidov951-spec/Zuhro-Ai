package service

import (
	"context"

	"zuhroai/internal/domain"
)

type ProductService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return nil, nil
}
