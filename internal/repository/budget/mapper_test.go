package budget

import (
	"testing"
	"time"
)

func TestBudgetDBModelRoundTrip(t *testing.T) {
	now := time.Now()
	db := BudgetDB{ID: 7, UserID: 2, Amount: 300, CurrencyID: 1, Description: "x", CreatedAt: now, UpdatedAt: now, ClosedAt: &now, PeriodStart: now, PeriodEnd: now}
	m := BudgetDBToModel(db)
	db2 := BudgetModelToDB(m)
	if db2.ID != db.ID || db2.UserID != db.UserID || db2.Amount != db.Amount || db2.CurrencyID != db.CurrencyID || db2.Description != db.Description {
		t.Fatalf("roundtrip mismatch: %+v vs %+v", db, db2)
	}
}
