package usecase

import "context"

type Notifier interface {
	Confirm(ctx context.Context, target string, id string) error
}
