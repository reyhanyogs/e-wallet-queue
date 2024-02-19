package domain

import (
	"context"

	"github.com/hibiken/asynq"
)

type AccountService interface {
	GenerateMutation() (string, func(ctx context.Context, task *asynq.Task) error)
}
