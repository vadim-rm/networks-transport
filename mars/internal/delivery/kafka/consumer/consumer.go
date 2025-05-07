package consumer

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type Handler func(ctx context.Context, messageValue []byte) error

type Consumer struct {
	reader  *kafka.Reader
	handler Handler
}

func NewConsumer(reader *kafka.Reader, handler Handler) *Consumer {
	return &Consumer{reader: reader, handler: handler}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message: %s", err.Error())
			continue
		}
		err = c.handler(ctx, m.Value)
		if err != nil {
			log.Printf("error handling message: %s", err.Error())
		}
	}
}
