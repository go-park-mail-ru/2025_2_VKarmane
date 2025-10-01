package operation

import (
	"testing"
	"time"
)

func TestOperationDBModelRoundTrip(t *testing.T) {
	now := time.Now()
	db := OperationDB{ID: 9, AccountID: 1, CategoryID: 2, Type: "expense", Status: "done", Description: "d", ReceiptURL: "u", Name: "n", Sum: 10.5, CurrencyID: 1, CreatedAt: now}
	m := OperationDBToModel(db)
	db2 := OperationModelToDB(m)
	if db2.ID != db.ID || db2.AccountID != db.AccountID || db2.CategoryID != db.CategoryID || db2.Type != db.Type || db2.Status != db.Status || db2.Sum != db.Sum || db2.CurrencyID != db.CurrencyID {
		t.Fatalf("roundtrip mismatch: %+v vs %+v", db, db2)
	}
}
