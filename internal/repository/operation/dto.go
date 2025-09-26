package operation

import "time"

type OperationDB struct {
	ID          int       `db:"operation_id"`
	AccountID   int       `db:"account_id"`
	CategoryID  int       `db:"category_id"`
	Type        string    `db:"type"`
	Status      string    `db:"status"`
	Description string    `db:"description"`
	ReceiptURL  string    `db:"receipt_url"`
	Name        string    `db:"name"`
	Sum         float64   `db:"sum"`
	CurrencyID  int       `db:"currency_id"`
	CreatedAt   time.Time `db:"created_at"`
}
