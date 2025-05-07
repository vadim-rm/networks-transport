package repository

import (
	"context"
	"transport/internal/domain"
)

type DataLink interface {
	Send(ctx context.Context, segment domain.Segment) error
}
