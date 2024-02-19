package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet-queue/domain"
)

type userRepository struct {
	db *goqu.Database
}

func NewUser(conn *sql.DB) domain.UserRepository {
	return &userRepository{
		db: goqu.New("default", conn),
	}
}

func (r *userRepository) FindEmailByID(ctx context.Context, id int64) (email string, err error) {
	dataset := r.db.From("users").Select("email").Where(goqu.Ex{
		"id": id,
	})
	_, err = dataset.ScanValContext(ctx, &email)
	return
}
