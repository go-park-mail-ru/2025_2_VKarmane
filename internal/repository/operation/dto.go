package operation

import "time"

type OperationDB struct {
	ID            int       `db:"operation_id"`
	AccountFromID *int      `db:"account_from_id"`
	AccountToID   *int      `db:"account_to_id"`
	CategoryID    *int      `db:"category_id"`
	CurrencyID    *int      `db:"currency_id"`
	Status        string    `db:"operation_status"`
	Type          string    `db:"operation_type"`
	Name          string    `db:"operation_name"`
	Description   string    `db:"operation_description"`
	ReceiptURL    string    `db:"receipt_url"`
	Sum           float64   `db:"sum"`
	CreatedAt     time.Time `db:"created_at"`
}
