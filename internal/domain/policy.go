package domain

import (
	"context"
)

type UniqueEmailPolicy interface {
	IsUnique(context.Context, Email) (bool, error)
}
