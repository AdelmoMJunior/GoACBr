# GoACBr API

GoACBr API is a high-performance, multi-tenant Go-based REST API that wraps the `ACBrLibNFe` C-library (using cgo). It is designed to handle electronic fiscal documents (NFe/NFCe) in Brazil, supporting the 2027 Tax Reform (Reforma Tributária - IS, IBS, CBS) out-of-the-box.

## Features

- **Multi-Tenant Architecture**: Supports multiple CNPJs simultaneously with thread-safe `ACBrLib` handle management.
- **PostgreSQL Database**: Normalized schema with over 15 tables to persist documents, items, transports, payments, and events.
- **Reforma Tributária Ready**: Full schema support for new taxes (`tot_is_v_is`, `tot_ibs_cbs_v_bc`, etc.) as required by ACBr `.ini` configurations.
- **Backblaze B2 Storage**: XML and PDF files are stored directly in B2 (S3-compatible) to keep the database lean.
- **Redis & In-Memory Fallback**: Uses Redis for session management, rate limiting, and rate-limit fallbacks. If Redis is down, it seamlessly falls back to PostgreSQL.
- **DFe Distribution Worker**: Background worker that queries SEFAZ for new documents, respecting the 1-hour cooldown per CNPJ.
- **Security**:
  - JWT Authentication with Refresh Tokens.
  - AES-256-GCM encryption for digital certificates (PFX) and SMTP passwords in the database.
  - Expiry checks before operations.

## Requirements

- **Go 1.24+**
- **Docker & Docker Compose** (for easy deployment of DB, Redis, and API)
- **ACBrLibNFe** shared library (`libacbrnfe64.so` for Linux) and its schemas.
- **Backblaze B2 Account** (for file storage)

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/AdelmoMJunior/GoACBr.git
   cd GoACBr
   ```

2. **Environment Variables:**
   Copy `.env.example` to `.env` and fill in the required keys, especially `JWT_SECRET`, `ENCRYPTION_KEY` (32 hex bytes), and `B2` credentials.
   ```bash
   cp .env.example .env
   ```

3. **Provide ACBrLib files:**
   Place the ACBr `libacbrnfe64.so` inside the `./lib` directory. Ensure that the schemas are correctly located in `./data/Schemas`.

4. **Run with Docker Compose:**
   ```bash
   make docker-up
   ```
   This will start PostgreSQL, Redis, and the Go API.

5. **Run Migrations:**
   ```bash
   make migrate-up
   ```

## Development

Use the provided `Makefile` for common tasks:
- `make build` - Builds the binary.
- `make run` - Runs the server locally.
- `make test` - Runs unit tests.
- `make lint` - Runs golangci-lint.
- `make swagger` - Generates Swagger documentation.

## Architecture Highlights

The project follows Clean Architecture principles:
- `cmd/api`: Entry point, server setup, dependency injection.
- `internal/domain`: Core domain models (User, Company, Invoice, etc.).
- `internal/config`: Configuration loader.
- `internal/crypto`: AES encryption and hashing utilities.
- `internal/storage`: B2 integration.
- `internal/repository`: PostgreSQL and Redis data access.
- `internal/service`: Core business logic and `ACBrLib` CGO wrapper logic.
- `internal/handler`: HTTP handlers and routers.
- `pkg/`: Generic utilities (logging, validators, httputils).

## License
MIT
