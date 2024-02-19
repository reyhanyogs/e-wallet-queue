package domain

import "context"

type UserRepository interface {
	FindEmailByID(ctx context.Context, id int64) (email string, err error)
}
