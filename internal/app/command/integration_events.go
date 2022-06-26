package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

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
	port string
}

func NewSendEmail(port string) SendEmail {
	return SendEmail{
		port: port,
	}
}

func (h SendEmail) Handle(ctx context.Context, e SendEmailCommand) error {
	fmt.Println("===> faking send email to", e.Email)
	go func() {
		time.Sleep(time.Second)
		fmt.Println("===> faking user confirmation")
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/registrations/%s", h.port, e.ID))
		if err != nil {
			fmt.Println("ERROR while calling confirmation:", err)
		}
		if resp.StatusCode != http.StatusOK {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("ERROR while reading body:", err)
			}
			fmt.Println("ERROR: response not OK\n", string(body))
		}
	}()
	return nil
}
