package infra

import (
	"context"
	"log"
	"time"

	"github.com/quintans/go-clean-ddd/internal/domain/usecase/command"
	"github.com/quintans/toolkit/latch"
)

func StartOutboxScheduler(ctx context.Context, lock *latch.CountDownLatch, heartbeat time.Duration, outboxUC command.FlushOutbox) {
	lock.Add(1)
	go func() {
		defer lock.Done()
		ticker := time.NewTicker(heartbeat)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := outboxUC.Handle(ctx)
				if err != nil {
					log.Printf("ERROR: failed to execute flush outbox: %s\n", err)
				}
			}
		}
	}()
}
