# Dashboard Pay2U Backend

Go backend service for the Pay2U Dashboard.

## Tech Stack

- **Go** (Fiber v2 + GORM + PostgreSQL + Redis)
- **JWT** authentication with RBAC authorization
- **Clean Architecture** (Delivery → UseCase → Repository → Entity)

## Getting Started

```bash
# 1. Copy config
cp config.yml.example config.yml

# 2. Edit config.yml with your database credentials

# 3. Run migrations
migrate -path db/migrations -database "postgres://user:pass@localhost:5432/dashboard_pay2u?sslmode=disable" up

# 4. Start the server
go run cmd/api/main.go
```

## Project Structure

```
cmd/api/main.go                        # Application entrypoint
internal/
├── common/constants/                  # Application-wide constants
├── common/hashing/                    # Bcrypt password hashing
├── config/                            # Fiber, GORM, Viper, Redis config
├── delivery/http/controller/          # HTTP handlers
├── delivery/http/middleware/          # Auth, RBAC, RequestID, Response middlewares
├── delivery/http/route/               # Route registration
├── entity/                            # GORM model entities
├── exception/                         # Error handler helpers
├── helpers/                           # Utility functions
├── logger/                            # Logrus context logger
├── model/converter/                   # Entity → DTO converters
├── model/payload/                     # Request/Response DTOs
├── repository/                        # Data access layer
└── usecase/                           # Business logic layer
db/migrations/                         # SQL migration files
```
