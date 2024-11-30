# Backend API with Fiber and SQLC

Backend API built with Go Fiber framework and SQLC.

## System Requirements

- Go 1.23.3 or higher
- PostgreSQL
- Air (for hot reload in development)
- SQLC
- Goose (for database migrations)

## Installation

1. Clone repository:
```bash
git clone <repository-url>
```

2. Install required packages:
```bash
go mod download
```

3. Create `.env` file and copy content from `.env.example`:
```bash
cp .env.example .env
```

4. Create database and run migrations:
```bash
goose -dir sqlc/migrations postgres "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

5. Generate SQLC code:
```bash
sqlc generate
```

6. Generate Swagger docs:
```bash
cd cmd/api
swag init -g main.go --parseDependency --parseInternal -o docs
```

## Running the Application

### Development Environment
```bash
set GO_ENV=development && air
```

### Production Environment
```bash
go build -o app ./cmd/api/main.go
```

## API Documentation

After running the application, you can access the Swagger UI at:
```
http://localhost:8386/swagger/
```

## Project Structure
```
.
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── db/              # SQLC generated code
│   ├── interfaces/      # Interface definitions
│   ├── models/          # Business models
│   ├── pkg/             # Shared packages
│   └── repository/      # Repository implementations
├── sqlc/
│   ├── migrations/      # Database migrations
│   ├── queries/         # SQLC queries
│   └── sqlc.yaml        # SQLC configuration
└── docs/               # API documentation
```

## Key Features

- JWT Authentication
- Role-based access control (RBAC)
- RESTful API endpoints
- Swagger documentation
- Hot reload in development
- Database migrations with Goose
- Type-safe database queries with SQLC