# Backend API với Fiber và GORM

Backend API được xây dựng bằng Go Fiber framework và GORM ORM.

## Yêu cầu hệ thống

- Go 1.23.3 trở lên
- PostgreSQL
- Air (cho hot reload trong môi trường development)

## Cài đặt

1. Clone repository:

```bash
git clone 
```

2. Cài đặt các package cần thiết:

```bash
go mod download
```

3. Tạo file `.env` và copy nội dung từ `.env.example` vào:

```bash
cp .env.example .env
```
cd cmd/api
swag init -g main.go --output ../../docs

## Chạy ứng dụng
### Môi trường Development
```bash
set GO_ENV=development && air
```

### Môi trường Production

```bash
go build -o app ./cmd/api/main.go
```
