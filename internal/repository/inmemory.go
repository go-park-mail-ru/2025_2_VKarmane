package repository

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/budget"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/operation"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

type Store struct {
	Users        []dto.UserDB
	Accounts     []account.AccountDB
	UserAccounts []account.UserAccountDB
	Operations   []operation.OperationDB
	Categories   []dto.CategoryDB
	Currency     []dto.CurrencyDB
	Budget       []budget.BudgetDB

	AccountRepo   *account.Repository
	BudgetRepo    *budget.Repository
	OperationRepo *operation.Repository
}

func NewStore() (*Store, error) {
	now := time.Now()

	// Пароль для всех тестовых пользователей: Test123
	testPassword := "Test123"

	passwordHash1, err := utils.HashPassword(testPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password for user 1: %w", err)
	}
	passwordHash2, err := utils.HashPassword(testPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password for user 2: %w", err)
	}

	users := []dto.UserDB{
		{ID: 1, FirstName: "Vlad", LastName: "Sigma", Email: "vlad@example.com", Login: "hello", Password: passwordHash1, CreatedAt: now, UpdatedAt: now},
		{ID: 2, FirstName: "Nikita", LastName: "Go", Email: "nikita@example.com", Login: "goodbye", Password: passwordHash2, CreatedAt: now, UpdatedAt: now},
	}
	currencies := []dto.CurrencyDB{
		{ID: 1, Code: "USD", Name: "US Dollar", CreatedAt: now},
		{ID: 2, Code: "RUB", Name: "Ruble", CreatedAt: now},
	}
	accounts := []account.AccountDB{
		{ID: 1, Balance: 100, Type: "card", CurrencyID: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Balance: 500, Type: "cash", CurrencyID: 1, CreatedAt: now, UpdatedAt: now},
	}
	userAccounts := []account.UserAccountDB{
		{ID: 1, UserID: 1, AccountID: 1, CreatedAt: now, UpdatedAt: now},
		{ID: 2, UserID: 1, AccountID: 2, CreatedAt: now, UpdatedAt: now},
		{ID: 3, UserID: 2, AccountID: 2, CreatedAt: now, UpdatedAt: now},
	}
	categories := []dto.CategoryDB{
		{ID: 1, UserID: 1, Name: "food", Description: "", CreatedAt: now, UpdatedAt: now},
		{ID: 2, UserID: 1, Name: "salary", Description: "income", CreatedAt: now, UpdatedAt: now},
	}
	operations := []operation.OperationDB{
		{ID: 1, AccountID: 1, CategoryID: 1, Type: "expense", Name: "Restaurant", Sum: 80, CurrencyID: 1, CreatedAt: now},
		{ID: 2, AccountID: 1, CategoryID: 1, Type: "expense", Name: "Vkusno i tochka", Sum: 30, CurrencyID: 1, CreatedAt: now},
		{ID: 3, AccountID: 1, CategoryID: 1, Type: "income", Name: "Salary", Sum: 70, CurrencyID: 1, CreatedAt: now},
	}
	budgets := []budget.BudgetDB{
		{ID: 1, UserID: 1, Amount: 100, CurrencyID: 1, Description: "September food", CreatedAt: now, UpdatedAt: now, PeriodStart: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC), PeriodEnd: time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC)},
		{ID: 2, UserID: 1, Amount: 500, CurrencyID: 1, Description: "September relax", CreatedAt: now, UpdatedAt: now, PeriodStart: time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC), PeriodEnd: time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC)},
	}
	store := &Store{
		Users:        users,
		Currency:     currencies,
		Accounts:     accounts,
		UserAccounts: userAccounts,
		Categories:   categories,
		Operations:   operations,
		Budget:       budgets,
	}

	clock := clock.RealClock{}
	store.AccountRepo = account.NewRepository(accounts, userAccounts, clock)
	store.BudgetRepo = budget.NewRepository(budgets, clock)
	store.OperationRepo = operation.NewRepository(operations, clock)

	return store, nil
}
