package service

import (
	"context"
	"transport/internal/domain"
)

type Transport interface {
	SendMessage(ctx context.Context, message domain.Message) error
}
