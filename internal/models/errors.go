package models

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrBudgetNotFound  = errors.New("budget not found")
	ErrUserNotFound    = errors.New("user not found")
)
