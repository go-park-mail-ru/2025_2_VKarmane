package operation

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type OperationDB struct {
	ID            int                    `db:"operation_id"`
	AccountFromID *int                   `db:"account_from_id"`
	AccountToID   *int                   `db:"account_to_id"`
	CategoryID    *int                   `db:"category_id"`
	CategoryName  string                 `db:"category_name"`
	CurrencyID    *int                   `db:"currency_id"`
	Status        models.OperationStatus `db:"operation_status"`
	Type          models.OperationType   `db:"operation_type"`
	Name          string                 `db:"operation_name"`
	Description   string                 `db:"operation_description"`
	ReceiptURL    string                 `db:"receipt_url"`
	Sum           float64                `db:"sum"`
	CreatedAt     time.Time              `db:"created_at"`
	Date          time.Time              `db:"operation_date"` // Дата операции
}
