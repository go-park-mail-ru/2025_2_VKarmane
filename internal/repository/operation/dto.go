package operation

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type OperationDB struct {
	ID          int                    `db:"operation_id"`
	AccountID   int                    `db:"account_id"`
	CategoryID  int                    `db:"category_id"`
	Type        models.OperationType   `db:"type"`
	Status      models.OperationStatus `db:"status"`
	Description string                 `db:"description"`
	ReceiptURL  string                 `db:"receipt_url"`
	Name        string                 `db:"name"`
	Sum         float64                `db:"sum"`
	CurrencyID  int                    `db:"currency_id"`
	CreatedAt   time.Time              `db:"created_at"`
}
