package service

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/reyhanyogs/e-wallet-queue/domain"
	"github.com/reyhanyogs/e-wallet-queue/internal/config"
)

type accountService struct {
	config *config.Config
}

func NewAccount(config *config.Config) domain.AccountService {
	return &accountService{
		config: config,
	}
}

func (s *accountService) GenerateMutation() (string, func(ctx context.Context, task *asynq.Task) error) {
	return "generate:mutation", func(ctx context.Context, task *asynq.Task) error {
		log.Println("generate mutation execute")
		return nil
	}
}
