package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	Price       float64    `json:"price" db:"price"`
	Stock       int        `json:"stock" db:"stock"`
	Category    string     `json:"category" db:"category"`
	ImageURL    string     `json:"image_url" db:"image_url"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedBy   uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"required"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" binding:"omitempty"`
	Description *string  `json:"description" binding:"omitempty"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
	Category    *string  `json:"category" binding:"omitempty"`
	ImageURL    *string  `json:"image_url" binding:"omitempty,url"`
	IsActive    *bool    `json:"is_active" binding:"omitempty"`
}
