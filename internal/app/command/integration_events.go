package command

import (
	"context"

	"github.com/quintans/go-clean-ddd/internal/app"
	"github.com/quintans/go-clean-ddd/internal/domain"
)

type SendEmailHandler interface {
	Handle(context.Context, SendEmailCommand) error
}

type SendEmailCommand struct {
	ID    string
	Email domain.Email
}

type SendEmail struct {
	sender  app.EmailSender
	rootUrl string
}

func NewSendEmail(rootUrl string, sender app.EmailSender) SendEmail {
	return SendEmail{
		sender:  sender,
		rootUrl: rootUrl,
	}
}

func (h SendEmail) Handle(ctx context.Context, e SendEmailCommand) error {
	return h.sender.Send(ctx, e.Email, h.rootUrl+e.ID)
}
