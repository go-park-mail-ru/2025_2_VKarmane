package models

import "encoding/json"

//go:generate easyjson -all message.go

const (
	TRANSACTIONS string = "transactions"
	CATEGORIES   string = "categories"
)

type KafkaMessageWrapper struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
