package fakeemail

import (
	"context"

	"github.com/quintans/go-clean-ddd/fake"
	"github.com/quintans/go-clean-ddd/internal/domain"
)

type Client struct {
	client fake.EmailClient
}

func NewClient(client fake.EmailClient) Client {
	return Client{
		client: client,
	}
}

func (f Client) Send(ctx context.Context, destination domain.Email, body string) error {
	f.client.Send(destination.String(), body)

	return nil
}
