package kafka

import (
	"context"
	"encoding/json"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaWriterWrapper struct {
	writer *kafka.Writer
}

type KafkaMessageWrapper struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewKafkaWriterWrapper(w *kafka.Writer) *KafkaWriterWrapper {
	return &KafkaWriterWrapper{writer: w}
}

func (k *KafkaWriterWrapper) WriteMessages(ctx context.Context, msgs ...KafkaMessage) error {
	kafkaMsgs := make([]kafka.Message, len(msgs))
	for i, m := range msgs {
		wrapper := KafkaMessageWrapper{
			Type:    m.Type,
			Payload: m.Payload,
		}

		data, _ := json.Marshal(wrapper)
		kafkaMsgs[i] = kafka.Message{Value: data}
	}
	return k.writer.WriteMessages(ctx, kafkaMsgs...)
}
