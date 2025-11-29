package models

type Category struct {
	CategoryID           int    `json:"category_id"`
	CategoryName         string `json:"category_name"`
	CategoryLogoHashedID string `json:"category_logo_hashed_id"`
	CategoryLogo         string `json:"category_logo"`
	Action               string `json:"action"`
}
