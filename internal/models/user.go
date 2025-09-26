package models

import "time"

type User struct {
	ID        int
	Name      string
	Surname   string
	Email     string
	Login     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
