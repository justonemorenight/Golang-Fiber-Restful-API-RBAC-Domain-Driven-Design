# Backend API with Fiber and GORM

Backend API is built with Go Fiber framework and GORM ORM.

## System requirements

- Go 1.23.3 or higher
- PostgreSQL
- Air (for hot reload in development environment)

## Installation

1. Clone repository:

```bash
git clone 
```

2. Install necessary packages:

```bash
go mod download
```

3. Create file `.env` and copy content from `.env.example`:

```bash
cp .env.example .env
```

4. Generate Swagger docs:

```bash
cd cmd/api
swag init -g main.go --parseDependency --parseInternal -o docs
```

## Run application

### Development environment

```bash
set GO_ENV=development && air
```

### Production environment

```bash
go build -o app ./cmd/api/main.go
```
