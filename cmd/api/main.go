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

	_ "backend-fiber/cmd/api/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// @title Backend Fiber API
// @version 1.0
// @description This is a sample server for Backend Fiber.
// @host localhost:8386
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize JWT config
	auth.InitJWTConfig(cfg)

	// Create connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Connect to database for migrations
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database for migrations:", err)
	}
	defer sqlDB.Close()

	// Always run migrations to update the latest schema
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal("Failed to set dialect:", err)
	}

	if err := goose.Up(sqlDB, "sqlc/migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Log current migration version
	current, err := goose.GetDBVersion(sqlDB)
	if err != nil {
		log.Fatal("Failed to get current migration version:", err)
	}
	log.Printf("Current migration version: %d", current)

	// Connect to database with pgxpool for the application
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbpool.Close()

	// Initialize queries from sqlc
	queries := sqlcdb.New(dbpool)

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(queries)
	refreshTokenRepo := repository.NewRefreshTokenRepository(queries, dbpool)
	userService := user.NewService(userRepo, refreshTokenRepo, queries)
	rbacService := rbac.NewService(queries)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService)
	roleHandler := handlers.NewRoleHandler(rbacService)

	// Initialize RBAC middleware
	rbacMiddleware := middleware.NewRBACMiddleware(rbacService)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public routes
	v1.Post("/register", userHandler.CreateUser)
	auth := v1.Group("/auth")
	auth.Post("/login", authHandler.Login)

	// Protected routes with RBAC
	users := v1.Group("/users")
	users.Use(middleware.Protected())

	// Routes for admin
	users.Get("/", rbacMiddleware.RequirePermission("users.list"), userHandler.GetUsers)
	users.Post("/", rbacMiddleware.RequirePermission("users.create"), userHandler.CreateUser)

	// Route for both admin and member
	users.Get("/:id", rbacMiddleware.RequirePermission("users.read_self"), userHandler.GetUserByID)

	// Role routes with RBAC
	roles := v1.Group("/roles")
	roles.Use(middleware.Protected())

	// Basic CRUD operations
	roles.Post("/", rbacMiddleware.RequirePermission("roles.create"), roleHandler.CreateRole)
	roles.Get("/", rbacMiddleware.RequirePermission("roles.list"), roleHandler.ListRoles)
	roles.Get("/:id", rbacMiddleware.RequirePermission("roles.read"), roleHandler.GetRole)
	roles.Put("/:id", rbacMiddleware.RequirePermission("roles.update"), roleHandler.UpdateRole)
	roles.Delete("/:id", rbacMiddleware.RequirePermission("roles.delete"), roleHandler.DeleteRole)

	// Permission management
	roles.Post("/:id/permissions", rbacMiddleware.RequirePermission("roles.assign_permission"), roleHandler.AssignPermission)
	roles.Delete("/:id/permissions/:permissionId", rbacMiddleware.RequirePermission("roles.remove_permission"), roleHandler.RemovePermission)
	roles.Get("/:id/permissions", rbacMiddleware.RequirePermission("roles.view_permissions"), roleHandler.ListRolePermissions)

	// User management
	roles.Post("/:id/users", rbacMiddleware.RequirePermission("roles.assign_user"), roleHandler.AssignUser)
	roles.Delete("/:id/users/:userId", rbacMiddleware.RequirePermission("roles.remove_user"), roleHandler.RemoveUser)
	roles.Get("/:id/users", rbacMiddleware.RequirePermission("roles.view_users"), roleHandler.ListRoleUsers)

	// Add route for swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Add new route to group users protected by middleware auth
	users.Get("/profile", rbacMiddleware.RequirePermission("users.read_self"), userHandler.GetProfile)

	log.Fatal(app.Listen(":8386"))
}
