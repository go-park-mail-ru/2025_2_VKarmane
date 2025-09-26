package balance

// AccountAPI - DTO для API ответов с аккаунтами
type AccountAPI struct {
	ID         int     `json:"account_id"`
	Balance    float64 `json:"balance"`
	Type       string  `json:"type"`
	CurrencyID int     `json:"currency_id"`
}

// BalanceAPI - DTO для API ответов с балансом
type BalanceAPI struct {
	UserID   int          `json:"user_id"`
	Accounts []AccountAPI `json:"accounts"`
}
