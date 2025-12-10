package service

import (
	"context"

	"suitemedia/internal/models"
	"suitemedia/internal/repository"
	"suitemedia/pkg/redis"
)

type ProductService interface {
	List(ctx context.Context, params models.ListParams) ([]*models.Product, int64, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	Create(ctx context.Context, req models.CreateProductRequest) (*models.Product, error)
	Update(ctx context.Context, id string, req models.UpdateProductRequest) (*models.Product, error)
	Delete(ctx context.Context, id string) error
}

type productService struct {
	productRepo repository.ProductRepository
	redis       *redis.Client
}

func NewProductService(productRepo repository.ProductRepository, redis *redis.Client) ProductService {
	return &productService{
		productRepo: productRepo,
		redis:       redis,
	}
}

func (s *productService) List(ctx context.Context, params models.ListParams) ([]*models.Product, int64, error) {
	return s.productRepo.List(ctx, params)
}

func (s *productService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) Create(ctx context.Context, req models.CreateProductRequest) (*models.Product, error) {
	// Implementation placeholder
	return nil, nil
}

func (s *productService) Update(ctx context.Context, id string, req models.UpdateProductRequest) (*models.Product, error) {
	// Implementation placeholder
	return nil, nil
}

func (s *productService) Delete(ctx context.Context, id string) error {
	return s.productRepo.Delete(ctx, id)
}
