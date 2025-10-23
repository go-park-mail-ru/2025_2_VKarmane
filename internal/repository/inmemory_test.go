package repository

import "testing"

func TestNewStore_InitializesRepositories(t *testing.T) {
	store, err := NewStore()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.AccountRepo == nil || store.BudgetRepo == nil || store.OperationRepo == nil {
		t.Fatalf("expected repos to be initialized")
	}
	if len(store.Users) == 0 || len(store.Accounts) == 0 || len(store.UserAccounts) == 0 || len(store.Budget) == 0 {
		t.Fatalf("expected seeded data")
	}
}
