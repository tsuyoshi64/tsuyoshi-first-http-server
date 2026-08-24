# Chirpy HTTP Server

A small Go HTTP API for a microblog-style app. It supports user accounts, JWT authentication, refresh tokens, chirp CRUD, a dev reset endpoint, static file serving, and a Polka webhook for upgrading users to Chirpy Red.

## Features

- User signup, login, and credential update
- Argon2id password hashing
- 1-hour JWT access tokens
- 60-day database-backed refresh tokens with revocation
- Authenticated chirp creation and deletion
- Chirp listing, sorting, filtering by author, and single-chirp lookup
- 140-character chirp limit with simple profanity filtering
- Chirpy Red upgrade webhook using an API key
- Static file serving from `/app/` with visit metrics
- Dev-only reset endpoint guarded by `PLATFORM=dev`

## Requirements

- Go `1.26.5` or compatible
- PostgreSQL
- `goose` for migrations
- `sqlc` only when regenerating database code

## Dependencies

Runtime Go dependencies are listed in [go.mod](./go.mod):

- `github.com/lib/pq` for PostgreSQL
- `github.com/joho/godotenv` for `.env` loading
- `github.com/google/uuid` for UUIDs
- `github.com/alexedwards/argon2id` for password hashing
- `github.com/golang-jwt/jwt/v5` for JWTs

## Configuration

Create a `.env` file in the project root:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=replace-me
POLKA_KEY=replace-me
PLATFORM=dev
```

`PLATFORM=dev` enables `POST /admin/reset`. Use any other value outside local development.

You can generate a nice long random string on the command line like this:
```sh
openssl rand -base64 64
```

## Database

Export `DB_URL`, then run migrations with goose:

```sh
export DB_URL="postgres://user:password@localhost:5432/chirpy?sslmode=disable"
goose -dir sql/schema postgres "$DB_URL" up
```

SQL queries live in [sql/queries](./sql/queries). Generated sqlc code lives in [internal/database](./internal/database). Regenerate it after changing SQL:

```sh
sqlc generate
```

## Usage

Run the server:

```sh
go run .
```

The server listens on `http://localhost:8080`.

Run tests:

```sh
go test ./...
```

## API

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/api/healthz` | None | Health check |
| `POST` | `/api/users` | None | Create a user with `email` and `password` |
| `PUT` | `/api/users` | Bearer JWT | Update the authenticated user's email and password |
| `POST` | `/api/login` | None | Login and receive an access token plus refresh token |
| `POST` | `/api/refresh` | Bearer refresh token | Receive a new access token |
| `POST` | `/api/revoke` | Bearer refresh token | Revoke a refresh token |
| `POST` | `/api/chirps` | Bearer JWT | Create a chirp |
| `GET` | `/api/chirps` | None | List chirps; supports `author_id`, `sort=asc`, and `sort=desc` |
| `GET` | `/api/chirps/{chirpID}` | None | Fetch one chirp |
| `DELETE` | `/api/chirps/{chirpID}` | Bearer JWT | Delete your own chirp |
| `POST` | `/api/polka/webhooks` | `ApiKey` | Handle `user.upgraded` webhook events |
| `GET` | `/admin/metrics` | None | Show static file hit count |
| `POST` | `/admin/reset` | Dev only | Delete users and reset metrics |

Example signup and chirp:

```sh
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"me@example.com","password":"secret"}'

curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"body":"hello chirpy"}'
```

## Architecture

- [main.go](./main.go) loads config, connects to PostgreSQL, wires routes, and starts the HTTP server.
- [apiConfig.go](./apiConfig.go) contains route handlers, response models, admin metrics, and webhook handling.
- [helper.go](./helper.go) contains JSON response helpers and chirp body cleaning.
- [internal/auth](./internal/auth) contains password hashing, JWT, refresh token, bearer token, and API key helpers.
- [internal/database](./internal/database) contains sqlc-generated query methods and models.
- [sql/schema](./sql/schema) contains goose migrations.
- [sql/queries](./sql/queries) contains sqlc query definitions.
