package service

type AccountView struct {
	ID         int     `json:"account_id"`
	Balance    float64 `json:"balance"`
	Type       string  `json:"type"`
	CurrencyID int     `json:"currency_id"`
}

type BalanceView struct {
	UserID   int           `json:"account_id"`
	Accounts []AccountView `json:"accounts"`
}

func (s *Service) GetBalanceForUser(userID int) (BalanceView, error) {
	accounts := s.store.GetAccountsByUser(userID)

	av := make([]AccountView, 0, len(accounts))

	for _, account := range accounts {
		av = append(av, AccountView{
			ID:         account.ID,
			Balance:    account.Balance,
			Type:       account.Type,
			CurrencyID: account.CurrencyID,
		})
	}

	return BalanceView{
		UserID:   userID,
		Accounts: av,
	}, nil
}
