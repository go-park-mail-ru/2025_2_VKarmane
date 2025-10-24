package operation

import (
	"testing"
	"time"
)

func TestOperationDBModelRoundTrip(t *testing.T) {
	now := time.Now()
	accountFromID := 1
	categoryID := 2
	currencyID := 1
	db := OperationDB{
		ID:            9,
		AccountFromID: &accountFromID,
		CategoryID:    &categoryID,
		Type:          "expense",
		Status:        "finished",
		Description:   "d",
		ReceiptURL:    "u",
		Name:          "n",
		Sum:           10.5,
		CurrencyID:    &currencyID,
		CreatedAt:     now,
		Date:          now,
		CategoryName:  "Test Category",
	}
	m := OperationDBToModel(db)
	db2 := OperationModelToDB(m)
	if db2.ID != db.ID || *db2.AccountFromID != *db.AccountFromID || *db2.CategoryID != *db.CategoryID || db2.Type != db.Type || db2.Status != db.Status || db2.Sum != db.Sum || *db2.CurrencyID != *db.CurrencyID {
		t.Fatalf("roundtrip mismatch: %+v vs %+v", db, db2)
	}
}
