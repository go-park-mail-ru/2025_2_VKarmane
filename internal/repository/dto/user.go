package dto

import "time"

type UserDB struct {
	ID        int       `db:"user_id"`
	Name      string    `db:"name"`
	Surname   string    `db:"surname"`
	Email     string    `db:"email"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
