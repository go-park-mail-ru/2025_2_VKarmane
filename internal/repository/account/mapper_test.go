package account

import (
	"testing"
	"time"
)

func TestAccountDBModelRoundTrip(t *testing.T) {
	now := time.Now()
	db := AccountDB{ID: 4, Balance: 20.5, Type: "card", CurrencyID: 2, CreatedAt: now, UpdatedAt: now}
	m := AccountDBToModel(db)
	db2 := AccountModelToDB(m)
	if db2.ID != db.ID || db2.Balance != db.Balance || db2.Type != db.Type || db2.CurrencyID != db.CurrencyID {
		t.Fatalf("roundtrip mismatch: %+v vs %+v", db, db2)
	}
}

func TestUserAccountDBModelRoundTrip(t *testing.T) {
	now := time.Now()
	db := UserAccountDB{ID: 3, UserID: 1, AccountID: 2, CreatedAt: now, UpdatedAt: now}
	m := UserAccountDBToModel(db)
	db2 := UserAccountModelToDB(m)
	if db2.ID != db.ID || db2.UserID != db.UserID || db2.AccountID != db.AccountID {
		t.Fatalf("roundtrip mismatch: %+v vs %+v", db, db2)
	}
}
