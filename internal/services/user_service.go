package services

import (
	apperrors "backend-fiber/internal/errors"
	"backend-fiber/internal/models"
	"backend-fiber/internal/repository"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
	if user.Name == "" || user.Email == "" {
		return apperrors.ErrValidation
	}

	existingUser, err := s.repo.GetByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrDatabaseOperation
	}
	if existingUser != nil {
		return apperrors.NewAppError(
			fiber.StatusConflict,
			"Email already exists",
			"A user with this email address already exists",
		)
	}

	return s.repo.Create(user)
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}
