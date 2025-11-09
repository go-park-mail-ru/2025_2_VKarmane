package models

import "time"

type Currency struct {
	ID        int
	Code      string
	Name      string
	CreatedAt time.Time
}
