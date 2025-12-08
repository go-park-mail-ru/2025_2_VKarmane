package models

import "time"

type Category struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	LogoHashedID string    `json:"logo_hashed_id,omitempty"`
	LogoURL      string    `json:"logo_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

type CreateCategoryReportRequest struct {
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

type CategoryInReport struct {
	CategoryID     int     `json:"category_id"`
	CategoryName   string  `json:"category_name"`
	OperationCount int     `json:"operation_count"`
	TotalSum       float64 `json:"total_sum"`
}

type CategoryReport struct {
	Categories []CategoryInReport `json:"categoires"`
	Start      time.Time          `json:"start_period"`
	End        time.Time          `json:"end_period"`
}
