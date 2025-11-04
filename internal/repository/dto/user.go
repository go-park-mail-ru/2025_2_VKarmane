package dto

import "time"

type UserDB struct {
	ID           int       `db:"user_id"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Email        string    `db:"email"`
	Login        string    `db:"login"`
	Password     string    `db:"password"`
	Description  string    `db:"description"`
	LogoHashedID string    `db:"logo_hashed_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
