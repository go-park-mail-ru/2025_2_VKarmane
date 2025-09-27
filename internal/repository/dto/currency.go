package dto

import "time"

type CurrencyDB struct {
	ID        int       `db:"currency_id"`
	Code      string    `db:"code"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}
