package models

import "time"

type Category struct {
	ID           int
	UserID       int
	Name         string
	Description  string
	LogoHashedID string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateCategoryRequest struct {
	Name         string `json:"name" validate:"required,max=30"`
	Description  string `json:"description,omitempty" validate:"max=60"`
	LogoHashedID string `json:"logo_hashed_id,omitempty"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,max=30"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=60"`
	LogoHashedID *string `json:"logo_hashed_id,omitempty"`
}

type CategoryWithStats struct {
	Category
	OperationsCount int `json:"operations_count"`
}
