package infra

import (
	"context"
	"time"

	"github.com/quintans/go-clean-ddd/internal/infra/controller/outbox"
)

func StartOutboxRelay(ctx context.Context, heartbeat time.Duration, outboxController outbox.OutboxController) {
	go func() {
		ticker := time.NewTicker(heartbeat)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				outboxController.Flush(ctx)
			}
		}
	}()
}
