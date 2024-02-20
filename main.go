package main

import (
	"github.com/hibiken/asynq"
	"github.com/reyhanyogs/e-wallet-queue/internal/component"
	"github.com/reyhanyogs/e-wallet-queue/internal/config"
	"github.com/reyhanyogs/e-wallet-queue/internal/repository"
	"github.com/reyhanyogs/e-wallet-queue/internal/service"
)

func main() {
	config := config.Get()

	redisConnection := asynq.RedisClientOpt{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Pass,
	}

	dbConnection := component.GetDatabaseConn(config)

	accountRepository := repository.NewAccount(dbConnection)
	userRepository := repository.NewUser(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)

	emailService := service.NewEmail(config)
	accountService := service.NewAccount(config, accountRepository, transactionRepository, userRepository, emailService)

	worker := asynq.NewServer(redisConnection, asynq.Config{
		Concurrency: 4,
	})
	mux := asynq.NewServeMux()
	mux.HandleFunc(emailService.SendEmailQueue())
	mux.HandleFunc(accountService.GenerateMutation())

	component.Log.Info("Starting Queue Worker")
	if err := worker.Run(mux); err != nil {
		component.Log.Fatalf("Main(Run): err = %s", err.Error())
	}
}
