package models

import "encoding/json"

const (
	TRANSACTIONS string = "transactions"
	CATEGORIES   string = "categories"
)

type KafkaMessageWrapper struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
