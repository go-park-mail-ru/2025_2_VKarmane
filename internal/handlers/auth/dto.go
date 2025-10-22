package auth

import "time"

type UserAPI struct {
	ID        int       `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// type User struct {
// 	ID        int
// 	FirstName string
// 	LastName  string
// 	Email     string
// 	Login     string
// 	Password  string
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }
