package service

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/reyhanyogs/e-wallet-queue/domain"
	"github.com/reyhanyogs/e-wallet-queue/internal/component"
	"github.com/reyhanyogs/e-wallet-queue/internal/config"
)

type accountService struct {
	config                *config.Config
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
	userRepository        domain.UserRepository
	emailService          domain.EmailService
}

func NewAccount(
	config *config.Config,
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	userRepository domain.UserRepository,
	emailService domain.EmailService,
) domain.AccountService {
	return &accountService{
		config:                config,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
		emailService:          emailService,
	}
}

func (s *accountService) GenerateMutation() (string, func(ctx context.Context, task *asynq.Task) error) {
	return "generate:mutation", func(ctx context.Context, task *asynq.Task) error {
		account_ids, err := s.accountRepository.GetAllAccountId(ctx)
		if err != nil {
			component.Log.Errorf("GenerateMutation(GetAllAccountId): err = %s", err.Error())
			return err
		}

		for _, id := range account_ids {
			component.Log.Infof("GenerateMutation is running: ID = %d", id)
			earning, err := s.transactionRepository.GetWeeklyEarning(ctx, id)
			if err != nil {
				component.Log.Errorf("GenerateMutation(GetWeeklyEarning): ID = %d: err = %s", id, err.Error())
				return err
			}
			spends, err := s.transactionRepository.GetWeeklySpending(ctx, id)
			if err != nil {
				component.Log.Errorf("GenerateMutation(GetWeeklySpending): ID = %d: err = %s", id, err.Error())
				return err
			}

			email, err := s.userRepository.FindEmailByID(ctx, id)
			if err != nil {
				component.Log.Errorf("GenerateMutation(FindEmailByID): ID = %d: err = %s", id, err.Error())
				return err
			}

			body := fmt.Sprintf("Here's your weekly earning & spending\nEarning: %f\nSpending: %f", earning, spends)
			err = s.emailService.Send(email, "[E-wallet] Weekly Earning & Spending", body)
			if err != nil {
				component.Log.Errorf("GenerateMutation(Send): ID = %d: err = %s", id, err.Error())
				return err
			}
		}

		return nil
	}
}
