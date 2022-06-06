package outbox

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/domain/usecase/command"
)

type OutboxController struct {
	flushOutboxHandler command.FlushOutboxHandler
}

func NewOutboxController(flushOutboxHandler command.FlushOutboxHandler) OutboxController {
	return OutboxController{
		flushOutboxHandler: flushOutboxHandler,
	}
}

func (c OutboxController) Flush(context.Context) error {
	c.flushOutboxHandler.Handle(context)
}
