package balance

type AccountAPI struct {
	ID         int     `json:"id"`
	Balance    float64 `json:"balance"`
	Type       string  `json:"type"`
	CurrencyID int     `json:"currency_id"`
	CreatedAt  string  `json:"created_at,omitempty"`
	UpdatedAt  string  `json:"updated_at,omitempty"`
}

type BalanceAPI struct {
	UserID   int          `json:"user_id"`
	Accounts []AccountAPI `json:"accounts"`
}
