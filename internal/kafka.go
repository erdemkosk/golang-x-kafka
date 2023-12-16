package internal

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Writer[T any] struct {
	w *kafka.Writer
}

func NewWriter[T any](addr, topic string) (Writer[T], func() error) {
	w := &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return Writer[T]{w: w}, w.Close
}

func (w *Writer[T]) WriteBatch(ctx context.Context, items ...T) error {
	messages := make([]kafka.Message, len(items))
	for i, item := range items {
		b, _ := json.Marshal(item) // using a naive approach for serialization
		messages[i] = kafka.Message{
			Value: b,
		}
	}
	return w.w.WriteMessages(ctx, messages...)
}
