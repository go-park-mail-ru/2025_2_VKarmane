package kafka

import "context"

type KafkaProducer interface {
	WriteMessages(ctx context.Context, msgs ...KafkaMessage) error
}

type KafkaMessage struct {
	Value []byte
}