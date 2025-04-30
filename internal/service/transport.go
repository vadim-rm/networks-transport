package service

import (
	"context"
	"transport/internal/domain"
)

type Transport interface {
	Run(ctx context.Context)
	ReceiveSegment(ctx context.Context, segment domain.Segment) error
	SendMessage(ctx context.Context, message domain.Message) error
}
