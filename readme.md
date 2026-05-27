# gobank

Simple banking API in Go with session-based auth, CSRF protection, and PostgreSQL.

## Tech

- Go 1.25
- [chi](https://github.com/go-chi/chi) for routing
- [pgx](https://github.com/jackc/pgx) + [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool) for Postgres
- [sqlc](https://sqlc.dev/) for type-safe queries
- [tern](https://github.com/jackc/tern) for migrations
- [scs](https://github.com/alexedwards/scs) for session management (Postgres-backed)
- [gorilla/csrf](https://github.com/gorilla/csrf) for CSRF protection
- [shopspring/decimal](https://github.com/shopspring/decimal) for precise decimal arithmetic
- bcrypt for password hashing

## Domain

- **Holder** — account titleholder, either `individual` (PF) or `company` (PJ).
- **Account** — where the balance lives. Belongs to one holder; one holder can own multiple accounts.
- **Transaction** — append-only record of every balance movement.

All money values use `NUMERIC(15,2)` mapped to `decimal.Decimal`. Balance-modifying operations run inside a transaction with `SELECT ... FOR UPDATE` to prevent races. Transfers lock both accounts in ascending UUID order to prevent deadlocks.

## Setup

Copy `.env.example` to `.env` and fill in the values. Then:

```bash
docker compose up -d
go run ./cmd/terndotenv
go run ./cmd/api
```

Server listens on `:3080`.

## Routes

All routes are prefixed with `/api/v1`.

### Public

- `POST /signup/individual` — create individual holder
- `POST /signup/company` — create company holder
- `POST /login` — authenticate, opens session
- `GET /csrftoken` — fetch CSRF token

### Authenticated

- `POST /logout` — destroy session
- `POST /account` — create account for the logged-in holder
- `GET /account/{id}/balance` — read balance
- `POST /account/{id}/deposit`
- `POST /account/{id}/withdrawal`
- `POST /account/transfer`
- `DELETE /account/{id}` — soft delete (only if balance is zero)

## CSRF

The middleware is enforced when `GOBANK_ENV` is not `dev`. In dev mode it stays disabled to ease testing in tools like Insomnia.