package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet-queue/domain"
)

type transactionRepository struct {
	db *goqu.Database
}

func NewTransaction(conn *sql.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: goqu.New("default", conn),
	}
}

func (r *transactionRepository) GetWeeklyEarning(ctx context.Context, id int64) (float64, error) {
	startDate := time.Now().AddDate(0, 0, -7)

	dataset := r.db.From("transactions").
		Select(goqu.SUM("amount").As("total")).
		Where(goqu.C("account_id").Eq(id),
			goqu.C("transaction_type").Eq("C"),
			goqu.C("transaction_datetime").Gte(startDate))

	var totalNullable sql.NullFloat64
	_, err := dataset.ScanValContext(ctx, &totalNullable)
	if err != nil {
		return 0, err
	}

	if totalNullable.Valid {
		return totalNullable.Float64, nil
	} else {
		return 0, nil
	}
}

func (r *transactionRepository) GetWeeklySpending(ctx context.Context, id int64) (float64, error) {
	startDate := time.Now().AddDate(0, 0, -7)

	dataset := r.db.From("transactions").
		Select(goqu.SUM("amount").As("total")).
		Where(goqu.C("account_id").Eq(id),
			goqu.C("transaction_type").Eq("D"),
			goqu.C("transaction_datetime").Gte(startDate))

	var totalNullable sql.NullFloat64
	_, err := dataset.ScanValContext(ctx, &totalNullable)
	if err != nil {
		return 0, err
	}

	if totalNullable.Valid {
		return totalNullable.Float64, nil
	} else {
		return 0, nil
	}
}
