package auth

import "time"

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Login        string    `json:"login"`
	Password     string    `json:"password,omitempty"`
	Description  string    `json:"description,omitempty"`
	LogoHashedID string    `json:"logo_hashed_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

type UpdateProfileRequest struct {
	UserID       int32 `json:"user_id"`
	FirstName    string `json:"first_name" validate:"required,max=50"`
	LastName     string `json:"last_name" validate:"required,max=50"`
	Email        string `json:"email" validate:"required,email"`
	LogoHashedID string `json:"logo_hashed_id,omitempty"`
}
