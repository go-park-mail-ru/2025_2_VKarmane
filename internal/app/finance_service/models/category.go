package models

import "time"

type Category struct {
	ID           int
	UserID       int
	Name         string
	Description  string
	LogoHashedID string
	LogoURL      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateCategoryRequest struct {
	UserID       int
	Name         string
	Description  string
	LogoHashedID string
}

type UpdateCategoryRequest struct {
	UserID       int
	CategoryID   int
	Name         *string
	Description  *string
	LogoHashedID *string
}

type CategoryWithStats struct {
	Category
	OperationsCount int
}
