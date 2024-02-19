package domain

import (
	"context"

	"github.com/hibiken/asynq"
)

type EmailService interface {
	SendEmailQueue() (string, func(ctx context.Context, task *asynq.Task) error)
	Send(to, subject, body string) error
}
