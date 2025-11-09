package dto

import "time"

type CategoryDB struct {
	ID           int       `db:"category_id"`
	UserID       int       `db:"user_id"`
	Name         string    `db:"category_name"`
	Description  *string   `db:"category_description"`
	LogoHashedID string    `db:"logo_hashed_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
