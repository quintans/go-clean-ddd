package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quintans/go-clean-ddd/internal/infra"
	"github.com/quintans/toolkit/latch"
)

func main() {
	lock := latch.NewCountDownLatch()

	cfg := infra.LoadEnvVars()
	ctx, cancel := context.WithCancel(context.Background())

	infra.Start(ctx, lock, cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	cancel()
	lock.WaitWithTimeout(3 * time.Second)
}
