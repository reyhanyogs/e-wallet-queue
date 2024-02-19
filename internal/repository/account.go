package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet-queue/domain"
)

type accountRepository struct {
	db *goqu.Database
}

// FindEmailByID implements domain.UserRepository.
func (*accountRepository) FindEmailByID(ctx context.Context, id int64) (email string, err error) {
	panic("unimplemented")
}

func NewAccount(conn *sql.DB) domain.AccountRepository {
	return &accountRepository{
		db: goqu.New("default", conn),
	}
}

func (r *accountRepository) GetAllAccountId(ctx context.Context) (account_ids []int64, err error) {
	dataset := r.db.From("accounts").Select("id").Order(goqu.I("id").Asc())
	err = dataset.ScanValsContext(ctx, &account_ids)
	if err != nil {
		return nil, err
	}
	return
}
