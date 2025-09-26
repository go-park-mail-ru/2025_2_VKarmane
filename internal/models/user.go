package models

import "time"

type User struct {
	ID        int
	Name      string
	Surname   string
	Email     string
	Login     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
