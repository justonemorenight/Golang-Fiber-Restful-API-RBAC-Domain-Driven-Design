package services

import (
	"backend-fiber/internal/db"
	apperrors "backend-fiber/internal/errors"
	"backend-fiber/internal/repository"
	"context"
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*db.User, error) {
	if name == "" || email == "" {
		return nil, apperrors.ErrValidation
	}

	// Kiểm tra email đã tồn tại
	_, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrDatabaseOperation
		}
		// Nếu là ErrNoRows thì tiếp tục tạo user mới
	} else {
		// Nếu không có lỗi tức là email đã tồn tại
		return nil, apperrors.NewAppError(
			fiber.StatusConflict,
			"Email already exists",
			"A user with this email address already exists",
		)
	}

	// Tạo user mới
	params := db.CreateUserParams{
		Name:     name,
		Email:    email,
		Password: password, // Trong thực tế nên hash password trước khi lưu
	}

	user, err := s.repo.Create(ctx, params)
	if err != nil {
		return nil, apperrors.ErrDatabaseOperation
	}

	return &user, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]db.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (*db.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseOperation
	}
	return &user, nil
}
