package user

import (
	"backend-fiber/internal/db"
	apperrors "backend-fiber/internal/errors"
	postgres "backend-fiber/internal/infrastructure/persistence/postgres"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"backend-fiber/internal/auth"
	"backend-fiber/internal/domain/user"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo         user.UserRepository
	refreshTokenRepo *postgres.RefreshTokenRepository
	queries          *db.Queries
}

func NewService(userRepo user.UserRepository, refreshTokenRepo *postgres.RefreshTokenRepository, queries *db.Queries) *Service {
	return &Service{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		queries:          queries,
	}
}

func (s *Service) CreateUser(ctx context.Context, name, email, password string) (*db.User, error) {
	if name == "" || email == "" {
		return nil, apperrors.ErrValidation
	}

	// Kiểm tra email đã tồn tại
	_, err := s.userRepo.GetByEmail(ctx, email)
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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	// Tạo user mới
	params := db.CreateUserParams{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	user, err := s.userRepo.Create(ctx, params)
	if err != nil {
		return nil, apperrors.ErrDatabaseOperation
	}

	// Lấy member role ID
	memberRole, err := s.queries.GetRoleByName(ctx, "member")
	if err != nil {
		return nil, apperrors.ErrDatabaseOperation
	}

	// Gán role member cho user mới
	err = s.queries.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: user.ID,
		RoleID: memberRole.ID,
	})
	if err != nil {
		return nil, apperrors.ErrDatabaseOperation
	}

	user.Password = ""
	return &user, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]db.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *Service) GetUserByID(ctx context.Context, id int32) (*db.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseOperation
	}
	return &user, nil
}

func (s *Service) Login(ctx context.Context, email, password string, ip string, userAgent string) (*db.User, string, string, error) {
	// Tìm user theo email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", "", apperrors.NewAppError(
			fiber.StatusUnauthorized,
			"Invalid credentials",
			"Email or password is incorrect",
		)
	}

	// Kiểm tra password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", apperrors.NewAppError(
			fiber.StatusUnauthorized,
			"Invalid credentials",
			"Email or password is incorrect",
		)
	}

	// Tạo access token
	accessToken, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", "", apperrors.NewAppError(
			fiber.StatusInternalServerError,
			"Token generation failed",
			"Could not generate access token",
		)
	}

	// Tạo refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, "", "", apperrors.NewAppError(
			fiber.StatusInternalServerError,
			"Token generation failed",
			"Could not generate refresh token",
		)
	}

	// Lưu refresh token vào database
	err = s.refreshTokenRepo.Create(ctx, user.ID, refreshToken,
		time.Now().Add(auth.GetJWTConfig().RefreshExp))
	if err != nil {
		return nil, "", "", apperrors.NewAppError(
			fiber.StatusInternalServerError,
			"Token storage failed",
			"Could not store refresh token",
		)
	}

	// Xóa password trước khi trả về
	user.Password = ""

	return &user, accessToken, refreshToken, nil
}

func (s *Service) ValidateRefreshToken(ctx context.Context, token string, ip string, userAgent string) error {
	// Lấy token từ database
	refreshToken, err := s.refreshTokenRepo.Get(ctx, token)
	if err != nil {
		return err
	}

	// Lấy thông tin usage
	usage, err := s.refreshTokenRepo.GetTokenUsage(ctx, refreshToken.ID)
	if err != nil {
		return err
	}

	// Kiểm tra các dấu hiệu bất thường
	if usage.UsageCount > 100 { // Quá nhiều lần refresh
		s.refreshTokenRepo.Delete(ctx, token)
		return errors.New("suspicious activity: high refresh count")
	}

	if time.Since(usage.LastUsedAt) < 1*time.Minute { // Refresh quá nhanh
		return errors.New("suspicious activity: too frequent refresh")
	}

	if usage.LastUsedIP != ip { // IP thay đổi
		log.Printf("Warning: Token used from different IP. Previous: %s, Current: %s",
			usage.LastUsedIP, ip)
	}

	// Update usage info
	err = s.refreshTokenRepo.UpdateTokenUsage(ctx, refreshToken.ID, ip, userAgent)
	if err != nil {
		log.Printf("Failed to update token usage: %v", err)
	}

	return nil
}
