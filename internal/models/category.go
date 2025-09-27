package models

import "time"

type Category struct {
	ID          int
	UserID      int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
