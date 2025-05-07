package repository

import (
	"context"
	"transport/internal/domain"
)

type Application interface {
	SendMessage(ctx context.Context, message domain.Message) error
}
