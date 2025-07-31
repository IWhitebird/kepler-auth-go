# Kepler Auth Go

High-performance authentication service built with Go, Gin, and GORM.

## Features

- **JWT Authentication** with proper middleware
- **User Management** with pagination and filtering
- **Email Service** integration
- **Swagger Documentation** at `/swagger/`
- **Clean Architecture** with proper separation of concerns
- **Docker Support** for easy deployment
- **PostgreSQL** database with GORM
- **Minimal Boilerplate** - straightforward, optimized code

## Quick Start

### Development

```bash
# Install dependencies
make deps

# Start development server with database
make dev

# Or run manually
make run
```

### Production

```bash
# Build and run with Docker
make docker-run

# Or build binary
make build
./bin/kepler-auth-go
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user profile
- `PATCH /api/auth/me` - Update current user profile
- `POST /api/auth/change-password` - Change password

### Users (Admin)
- `GET /api/users` - List users with pagination/filtering
- `GET /api/users/:id` - Get user by ID
- `PATCH /api/users/:id` - Update user (admin only)
- `DELETE /api/users/:id` - Delete user (admin only)

### Email
- `POST /api/email/send` - Send email

## Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=auth05
DB_USER=postgres
DB_PASSWORD=postgres

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=86400

# Server
PORT=8000
GIN_MODE=debug
```

## Documentation

- **Swagger UI**: `http://localhost:8000/swagger/`
- **Health Check**: `http://localhost:8000/health`

## Make Commands

```bash
make deps          # Install dependencies
make build         # Build binary
make run           # Run development server
make test          # Run tests
make swagger       # Generate Swagger docs
make docker-build  # Build Docker image
make docker-run    # Run with Docker Compose
make dev           # Start development environment
make clean         # Clean build artifacts
```

## Project Structure

```
kepler-auth-go/
├── cmd/main.go                    # Application entry point
├── internal/
│   ├── api/                       # HTTP server and routes
│   ├── config/                    # Configuration management
│   ├── database/                  # Database connection
│   ├── handlers/                  # HTTP handlers
│   ├── middleware/                # Authentication & CORS
│   ├── models/                    # Data models & DTOs
│   └── services/                  # Business logic
├── docs/                          # Swagger documentation
├── Dockerfile                     # Container definition
├── docker-compose.yml             # Development environment
└── Makefile                       # Build commands
```