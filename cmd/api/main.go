package main

// @title Fiber API
// @version 1.0
// @description API documentation cho ứng dụng Fiber
// @host localhost:8386
// @BasePath /api/v1
import (
	"backend-fiber/internal/config"
	"backend-fiber/internal/handlers"
	"backend-fiber/internal/middleware"
	"backend-fiber/internal/models"
	"backend-fiber/internal/repository"
	"backend-fiber/internal/services"
	"fmt"
	"log"

	_ "backend-fiber/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{})

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(db)
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
