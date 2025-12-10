package models

import "testing"

func TestProductModel(t *testing.T) {
	product := &Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Stock:       10,
		Category:    "test",
	}

	if product.Name != "Test Product" {
		t.Errorf("Expected name 'Test Product', got %s", product.Name)
	}
	if product.Price != 99.99 {
		t.Errorf("Expected price 99.99, got %f", product.Price)
	}
}

func TestCreateProductRequest(t *testing.T) {
	req := CreateProductRequest{
		Name:        "New Product",
		Description: "Description",
		Price:       199.99,
		Stock:       5,
		Category:    "electronics",
	}

	if req.Name != "New Product" {
		t.Errorf("Expected name 'New Product', got %s", req.Name)
	}
}
