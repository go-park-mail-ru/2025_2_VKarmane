package account

import "time"

type AccountDB struct {
	ID         int       `db:"account_id"`
	Balance    float64   `db:"balance"`
	Type       string    `db:"type"`
	CurrencyID int       `db:"currency_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type UserAccountDB struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	AccountID int       `db:"account_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
