package repository

import (
	"backend-fiber/internal/db"
	"context"
)

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *UserRepository) Create(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *UserRepository) GetAll(ctx context.Context) ([]db.User, error) {
	return r.queries.ListUsers(ctx)
}

func (r *UserRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.queries.GetUser(ctx, id)
}
