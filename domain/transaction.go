package domain

import "context"

type TransactionRepository interface {
	GetWeeklyEarning(ctx context.Context, id int64) (float64, error)
	GetWeeklySpending(ctx context.Context, id int64) (float64, error)
}
