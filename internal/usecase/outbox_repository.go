package usecase

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/entity"
)

type OutboxRepository interface {
	Save(ctx context.Context, ob entity.Outbox) error
	LockAndLoad(ctx context.Context) ([]entity.Outbox, error)
	Consume(ctx context.Context, outboxes []entity.Outbox) error
	Release(ctx context.Context, outboxes []entity.Outbox) error
}
