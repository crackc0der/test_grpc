package product

import (
	"context"
	"fmt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetProducts(ctx context.Context) ([]*Product, error) {
	products, err := s.repository.SelectAllProducts(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting products in service's method GetProducts: %v", err)
	}

	return products, nil
}

func (s *Service) GetProduct(ctx context.Context, id int) (*Product, error) {
	product, err := s.repository.SelectProductById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting product in service's method GetProduct: %v", err)
	}

	return product, nil
}

func (s *Service) AddProduct(ctx context.Context, product *Product) (*Product, error) {
	product, err := s.repository.InsertProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("error updating product in service's method UpdateProduct: %v", err)
	}

	return product, nil
}

func (s *Service) DeleteProduct(ctx context.Context, id int) error {
	err := s.repository.DeleteProductById(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting product in service's method DeleteProduct: %v", err)
	}

	return nil
}

func (s *Service) UpdateProduct(ctx context.Context, product *Product) (*Product, error) {
	product, err := s.repository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("error updating product in service's method UpdateProduct: %v", err)
	}

	return product, nil
}
