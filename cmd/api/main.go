package main

import (
	"backend-fiber/internal/config"
	sqlcdb "backend-fiber/internal/db"
	"backend-fiber/internal/handlers"
	"backend-fiber/internal/middleware"
	"backend-fiber/internal/repository"
	"backend-fiber/internal/services"
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

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Tạo connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Kết nối database cho goose migrations
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database for migrations:", err)
	}
	defer sqlDB.Close()

	// Kiểm tra version hiện tại của migrations
	current, err := goose.GetDBVersion(sqlDB)
	if err != nil {
		log.Fatal("Failed to get migration version:", err)
	}

	// Chỉ chạy migrations nếu chưa được áp dụng
	if current == 0 {
		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatal("Failed to set dialect:", err)
		}

		if err := goose.Up(sqlDB, "db/migrations"); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
	}

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
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	users := v1.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetUsers)
	users.Get("/:id", userHandler.GetUserByID)

	// Thêm route cho swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":8386"))
}
