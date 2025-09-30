package models

import "time"

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Login     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Login    string `json:"login" validate:"required,min=3,max=30,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
