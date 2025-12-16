package kafka

import "context"

type KafkaProducer interface {
	WriteMessages(ctx context.Context, msgs ...KafkaMessage) error
}

type KafkaMessage struct {
	Type    string
	Payload []byte
}
