package dto

import "time"

type CategoryDB struct {
	ID          int       `db:"category_id"`
	UserID      int       `db:"user_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
