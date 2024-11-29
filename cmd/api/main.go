package main

import (
	"backend-fiber/internal/application/rbac"
	user "backend-fiber/internal/application/user"
	"backend-fiber/internal/auth"
	sqlcdb "backend-fiber/internal/db"
	repository "backend-fiber/internal/infrastructure/persistence/postgres"
	"backend-fiber/internal/interfaces/http/handlers"
	"backend-fiber/internal/interfaces/http/middleware"
	"backend-fiber/internal/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "backend-fiber/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// @title Backend API
// @version 1.0
// @description Backend API with Fiber
// @host localhost:8386
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Khởi tạo JWT config
	auth.InitJWTConfig(cfg)

	// Tạo connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Kết nối database cho goose migrations
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database for migrations:", err)
	}
	defer sqlDB.Close()

	// Luôn chạy migration để cập nhật schema mới nhất
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("Failed to set dialect:", err)
	}

	if err := goose.Up(sqlDB, "sqlc/migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Log migration version hiện tại
	current, err := goose.GetDBVersion(sqlDB)
	if err != nil {
		log.Fatal("Failed to get current migration version:", err)
	}
	log.Printf("Current migration version: %d", current)

	// Kết nối database với pgxpool cho ứng dụng
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbpool.Close()

	// Khởi tạo queries từ sqlc
	queries := sqlcdb.New(dbpool)

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(queries)
	refreshTokenRepo := repository.NewRefreshTokenRepository(queries, dbpool)
	userService := user.NewService(userRepo, refreshTokenRepo, queries)
	userHandler := handlers.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService)

	// Public routes
	v1.Post("/register", userHandler.CreateUser)
	auth := v1.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// Khởi tạo RBAC service
	rbacService := rbac.NewService(queries)
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService)

	// Protected routes with RBAC
	users := v1.Group("/users")
	users.Use(middleware.Protected())

	// Routes cho admin
	users.Get("/", rbacMiddleware.RequirePermission("users.list"), userHandler.GetUsers)
	users.Post("/", rbacMiddleware.RequirePermission("users.create"), userHandler.CreateUser)

	// Route cho cả admin và member
	users.Get("/:id", rbacMiddleware.RequirePermission("users.read_self"), userHandler.GetUserByID)

	// Thêm route cho swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Thêm route mới vào group users đã được bảo vệ bởi middleware auth
	users.Get("/profile", rbacMiddleware.RequirePermission("users.read_self"), userHandler.GetProfile)

	log.Fatal(app.Listen(":8386"))
}
