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

# 2. Edit config.yml with your database credentials and ms_gateway.base_url

# 3. Run migrations
go run ./cmd/migrate up

# 4. Start the server
go run ./cmd/api/main.go
```

## Database migrations and seed data

Run these commands from the backend directory. The migration runner reads PostgreSQL credentials from `config.yml` and tracks Dashboard migrations in `dashboard_schema_migrations`, separately from the MS service. Install the `migrate` CLI from golang-migrate first:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

```bash
go run ./cmd/migrate up
```

Dashboard and MS must point to the same PostgreSQL database. Apply Dashboard migrations before the MS migration that adds the merchant-to-admin foreign key, then run the Dashboard seeder.

To roll back only the latest migration:

```bash
go run ./cmd/migrate down 1
```

If a failed first migration leaves the Dashboard tracker dirty, reset it to the no-migrations state and rerun:

```bash
go run ./cmd/migrate force -1
go run ./cmd/migrate up
```

The RBAC migration creates the default roles and permissions. To create a super-admin account, set its password and run the seeder:

```yaml
# config.yml
seed:
  admin:
    username: "admin"
    password: "set-a-strong-password"
    full_name: "System Administrator"
    email: "admin@example.com"
```

Then run:

```bash
go run ./cmd/seed
```

The password is required; username and full name have defaults, and email is optional. Running the seeder again does not replace an existing account's password or profile, and does not duplicate its role or permissions. Run it after migrations and from this directory.

Super Admins can list, create, update, and deactivate Admin accounts. Deactivation is soft so existing merchants retain their owner FK. Each admin and super admin has a server-generated MS key stored on the user record; the key is never returned to the frontend. The Dashboard backend sends the signed-in user's key to `ms-gateway-pay2u`. Both services must use the same PostgreSQL database so MS can validate the key against `users.ms_master_key`. Admins manage their own merchants and see transactions and analytics for those merchants; merchant records and credentials remain managed by MS.

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
