package models

type ESHitSource struct {
	ID                   int32   `json:"id"`
	AccountID            int32   `json:"account_id"`
	CategoryID           int32   `json:"category_id"`
	CategoryName         string  `json:"category_name"`
	Type                 string  `json:"type"`
	Description          string  `json:"description"`
	Name                 string  `json:"name"`
	CategoryLogoHashedId string  `json:"category_logo_hashed_id"`
	CategoryLogo         string  `json:"category_logo"`
	Sum                  float64 `json:"sum"`
	AccountType          string  `json:"account_type"`
	CurrencyId           int32   `json:"curerncy"`
	CreatedAt            string  `json:"created_at"`
	Date                 string  `json:"date"`
}

type ElasticsearchResponse struct {
	Hits struct {
		Hits []struct {
			Source ESHitSource `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
