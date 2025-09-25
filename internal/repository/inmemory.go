package repository

import (
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type Store struct {
	Users        []models.User
	Accounts     []models.Account
	UserAccounts []models.UserAccount
	Operations   []models.Operation
	Categories   []models.Category
	Currency     []models.Currency
	Budget       []models.Budget
}

func NewStore() *Store {
	now := time.Now()
	users := []models.User{
		{ID: 1, Name: "Vlad", Surname: "Sigma", Email: "example@example.com", Login: "hello", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Nikita", Surname: "Go", Email: "example@example.com", Login: "goodbye", CreatedAt: now, UpdatedAt: now},
	}
	currencies := []models.Currency{
		{ID: 1, Code: "USD", Name: "US Dollar", CreatedAt: now},
		{ID: 2, Code: "RUB", Name: "Ruble", CreatedAt: now},
	}
	accounts := []models.Account{
		{ID: 1, Balance: 100, Type: "card", CurrencyID: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Balance: 500, Type: "cash", CurrencyID: 1, CreatedAt: now, UpdatedAt: now},
	}
	userAccounts := []models.UserAccount{
		{ID: 1, UserID: 1, AccountID: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 2, UserID: 1, AccountID: 2, CreatedAt: now, UpdatedAt: now},
		{ID: 3, UserID: 2, AccountID: 2, CreatedAt: now, UpdatedAt: now},
	}
	categories := []models.Category{
		{ID: 1, UserID: 1, Name: "food", Description: "", CreatedAt: now, UpdatedAt: now},
		{ID: 2, UserID: 1, Name: "salary", Description: "income", CreatedAt: now, UpdatedAt: now},
	}
	operations := []models.Operation{
		{ID: 1, AccountID: 1, CategoryID: 1, Type: "expense", Name: "Restaurant", Sum: 80, CurrencyID: 1, CreatedAt: now},
		{ID: 2, AccountID: 1, CategoryID: 1, Type: "expense", Name: "Vkusno i tochka", Sum: 30, CurrencyID: 1, CreatedAt: now},
		{ID: 3, AccountID: 1, CategoryID: 1, Type: "income", Name: "Salary", Sum: 70, CurrencyID: 1, CreatedAt: now},
	}
	budgets := []models.Budget{
		{ID: 1, UserID: 1, Amount: 100, CurrencyID: 1, Description: "September food", CreatedAt: now, UpdatedAt: now, PeriodStart: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC), PeriodEnd: time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC)},
		{ID: 2, UserID: 1, Amount: 500, CurrencyID: 1, Description: "September relax", CreatedAt: now, UpdatedAt: now, PeriodStart: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC), PeriodEnd: time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC)},
	}
	return &Store{
		Users:        users,
		Currency:     currencies,
		Accounts:     accounts,
		UserAccounts: userAccounts,
		Categories:   categories,
		Operations:   operations,
		Budget:       budgets,
	}
}
