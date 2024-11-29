package user

import (
	"backend-fiber/internal/db"
	"context"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (db.User, error)
	Create(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetAll(ctx context.Context) ([]db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
}
