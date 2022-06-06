package usecase

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
)

type OutboxRepository interface {
	Save(ctx context.Context, ob entity.Outbox) error
	Consume(ctx context.Context, handler func([]entity.Outbox) error) error
}
