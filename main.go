package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/reyhanyogs/e-wallet-queue/internal/config"
	"github.com/reyhanyogs/e-wallet-queue/internal/service"
)

func main() {
	config := config.Get()

	redisConnection := asynq.RedisClientOpt{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Pass,
	}

	emailService := service.NewEmail(config)
	accountService := service.NewAccount(config)

	worker := asynq.NewServer(redisConnection, asynq.Config{
		Concurrency: 4,
	})
	mux := asynq.NewServeMux()
	mux.HandleFunc(emailService.SendEmailQueue())
	mux.HandleFunc(accountService.GenerateMutation())

	if err := worker.Run(mux); err != nil {
		log.Fatal(err)
	}
}
