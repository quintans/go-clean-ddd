package fakesub

import (
	"context"
	"encoding/json"

	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/fake"
	"github.com/quintans/go-clean-ddd/internal/app/command"
	"github.com/quintans/go-clean-ddd/internal/domain"
)

// RegistrationController handles the MQ about the creation of a new registration
// by sending an email of the registration
type RegistrationController struct {
	sendEmailHandler command.SendEmailHandler
}

func NewRegistrationController(sendEmailHandler command.SendEmailHandler) RegistrationController {
	return RegistrationController{
		sendEmailHandler: sendEmailHandler,
	}
}

func (h RegistrationController) Handle(ctx context.Context, e fake.MQEvent) error {
	var event RegistrationCreatedEvent
	err := json.Unmarshal(e.Payload, &event)
	if err != nil {
		return faults.Wrap(err)
	}

	email, err := domain.NewEmail(event.Email)
	if err != nil {
		return faults.Wrap(err)
	}

	return h.sendEmailHandler.Handle(ctx, command.SendEmailCommand{
		ID:    event.ID,
		Email: email,
	})
}

type RegistrationCreatedEvent struct {
	ID    string
	Email string
}
