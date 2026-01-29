# Go Architecture - DDD Clean Architecture

A modern Go REST API following Domain-Driven Design (DDD) principles with clean architecture.

## 🏗️ Architecture

```
go-architecture/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── product/                 # Product domain module
│   │   ├── domain/              # Business logic & entities
│   │   │   ├── product.go       # Product entity
│   │   │   ├── price.go         # Value object
│   │   │   └── repository.go    # Repository interface
│   │   ├── application/         # Use cases & DTOs
│   │   │   ├── dto.go           # Data transfer objects
│   │   │   ├── service.go       # Application services
│   │   │   ├── validator.go     # Input validation
│   │   │   └── mapper.go        # Domain <-> DTO mapping
│   │   └── infra/               # Infrastructure layer
│   │       ├── http/            # HTTP handlers
│   │       │   └── handler.go
│   │       └── postgres/        # Database implementation
│   │           └── repository.go
│   └── shared/                  # Shared infrastructure
│       ├── config/              # Configuration management
│       ├── errors/              # Error handling
│       ├── logger/              # Logging
│       └── middleware/          # HTTP middleware
│           ├── error_handler.go
│           ├── jwt.go           # JWT authentication
│           ├── logger.go        # Request logging
│           └── rate_limiter.go  # Rate limiting
├── migrations/                  # Database migrations
├── go.mod
└── .env.example
```

## ✨ Features

- **Clean Architecture**: Separation of concerns with clear boundaries
- **DDD Principles**: Rich domain model with entities, value objects, and repositories
- **Security**: JWT authentication, rate limiting, CORS
- **Validation**: Request validation using go-playground/validator
- **Error Handling**: Centralized error handling with custom error types
- **Logging**: Structured JSON logging
- **Database**: PostgreSQL with connection pooling
- **API Documentation**: RESTful API design

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

3. Create the database:

```sql
CREATE DATABASE goarch;
```

4. Run migrations:

```bash
psql -U postgres -d goarch -f migrations/001_create_products_table.sql
```

5. Install dependencies:

```bash
go mod download
```

6. Run the application:

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## 📚 API Endpoints


## 🔐 Auth — Login endpoint

- Endpoint: `POST /api/v1/login`
- Purpose: authenticate user and return JWT access token.
- Request (JSON):
```json
{
  "username": "admin",
  "password": "admin123"
}
```
- Response (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_at": "2026-01-30T16:56:59Z"
}
```
- Use: include header `Authorization: Bearer <access_token>` on protected endpoints (e.g., POST /api/v1/products).
- Notes:
  - Current implementation: demo in-memory credential check. Replace with real user store and bcrypt password checks for production.
  - JWT secret is configured via `JWT_SECRET` in `.env`.
  - Token expiration configured via `JWT_EXPIRATION` (hours).

### Example (curl)
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Example (Postman)
- Method: POST
- URL: http://localhost:8080/api/v1/login
- Body → raw → JSON:
```json
{
  "username": "admin",
  "password": "admin123"
}
```
- Response: copy `access_token` and use Authorization → Bearer Token for subsequent requests.

### Products

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/products` | No | List all products |
| GET | `/api/v1/products/:id` | No | Get product by ID |
| POST | `/api/v1/products` | Yes | Create product |
| PUT | `/api/v1/products/:id` | Yes | Update product |
| DELETE | `/api/v1/products/:id` | Yes | Delete product |

### Health Check

- `GET /health` - Health check endpoint

## 🔐 Authentication

Protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## 📝 Example Requests

### Create Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Laptop Pro",
    "description": "High-performance laptop",
    "price": 1299.99,
    "stock": 50,
    "category": "Electronics"
  }'
```

### Get All Products

```bash
curl http://localhost:8080/api/v1/products?category=Electronics&limit=10
```

### Update Product

```bash
curl -X PUT http://localhost:8080/api/v1/products/{id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Laptop Pro Max",
    "description": "Ultra high-performance laptop",
    "price": 1599.99,
    "stock": 30,
    "category": "Electronics"
  }'
```

## 🧪 Testing

Run tests:

```bash
go test ./...
```

## 🏛️ Architecture Principles

### Domain Layer
- Contains business logic and rules
- No dependencies on external layers
- Pure Go with no framework dependencies

### Application Layer
- Orchestrates domain objects
- Defines use cases
- Handles validation and DTOs

### Infrastructure Layer
- Implements interfaces defined in domain
- HTTP handlers, database repositories
- Framework-specific code

### Shared Layer
- Cross-cutting concerns
- Configuration, logging, error handling
- Middleware and utilities

## 📦 Dependencies

- **Fiber**: Fast HTTP framework
- **SQLX**: SQL extensions for Go
- **JWT**: JSON Web Token authentication
- **Validator**: Struct validation
- **PostgreSQL**: Database driver

## 🔒 Security Features

- JWT-based authentication
- Rate limiting (100 requests/minute per IP)
- CORS configuration
- Input validation
- SQL injection prevention (parameterized queries)
- Structured error responses (no sensitive data leakage)

## 🎯 Best Practices

- **SOLID principles** implemented
- **Dependency injection** for loose coupling
- **Interface-based design** for testability
- **Value objects** for domain concepts
- **Rich domain model** with encapsulated business logic
- **Graceful shutdown** for resource cleanup

## 📄 License

MIT
