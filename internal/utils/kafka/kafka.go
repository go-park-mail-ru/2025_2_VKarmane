package kafka

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaWriterWrapper struct {
	writer *kafka.Writer
}

func NewKafkaWriterWrapper(w *kafka.Writer) *KafkaWriterWrapper {
	return &KafkaWriterWrapper{writer: w}
}

func (k *KafkaWriterWrapper) WriteMessages(ctx context.Context, msgs ...KafkaMessage) error {
	kafkaMsgs := make([]kafka.Message, len(msgs))
	for i, m := range msgs {
		kafkaMsgs[i] = kafka.Message{Value: m.Value}
	}
	return k.writer.WriteMessages(ctx, kafkaMsgs...)
}