# Go Kickstart

A monorepo for a Go API and a Vite + React web app with shared TypeScript packages, managed with Turborepo and Bun workspaces.

## Repository layout

```
go-kickstart/
├── apps/api             # Go API (Fiber)
├── apps/web             # Vite + React frontend
├── packages/zod         # Shared Zod schemas
├── packages/openapi     # OpenAPI generation
├── packages/emails      # React Email templates
└── packages/*           # Other shared packages
```

## Prerequisites

- Go 1.24+
- Bun 1.2.13 (Node 22+)
- PostgreSQL 16+
- Redis 8+

## Quick start

```bash
bun install
cp apps/api/.env.example apps/api/.env
# Start PostgreSQL and Redis, then run migrations:
bun run api:migrate:up

# Start all apps
bun dev
```

Run the API only with `bun run api:run`.

## Common commands

```bash
# Monorepo (from root)
bun dev         # Start dev servers for all apps
bun dev:all    # Start dev servers for all apps and packages
bun run test       # Run tests for all apps and packages
bun run build
bun run lint
bun run typecheck

# API helpers (see apps/api/Makefile for migrate targets)
bun run api:run
bun run api:test
cd apps/api && make migrate-new NAME=add_table
cd apps/api && make migrate-up
cd apps/api && make migrate-down

# Contracts and emails
bun run openapi:generate    # Generate OpenAPI spec file from contracts
bun run emails:generate     # Generate email HTML templates

# UI components
bun run web:shadcn:add <component>
```

## API (apps/api)

- Clean layers: handlers -> services -> repositories -> models.
- Repositories are data access only; services implement business rules and validations.
- Use `ResourceRepository` / `ResourceService` / `ResourceHandler` for standard CRUD models.
- Entry points: `apps/api/cmd/api/main.go` (server) and `apps/api/cmd/seed/main.go` (seeder).
- Routes in `apps/api/internal/router/routes.go`; middleware order in `apps/api/internal/router/router.go`.
- Prefer `handler.Handle` / `handler.HandleNoContent` / `handler.HandleFile` for new endpoints.
- Request DTOs implement `validation.Validatable`; use `validation.BindAndValidate` or the handler wrappers.
- Use `utils.ParseUUIDParam` for `:id` params.
- Services return `errs.ErrorResponse`; wrap DB errors with `sqlerr.HandleError`. Handlers return errors and let `GlobalErrorHandler` format responses.
- Request IDs are set in middleware and injected into logs; use `middleware.GetLogger` in handlers.
- Context timeouts should use `server.Config.Server.ReadTimeout` / `WriteTimeout`.
- Auth uses short-lived JWT access tokens and long-lived refresh tokens. `middleware.Auth.RequireAuth` sets `user_id` in Fiber locals; sessions live in `auth_sessions`. Cookie config is under `AuthConfig`.
- Auth routes: `/api/v1/auth/register`, `/login`, `/google`, `/verify-email`, `/refresh`, `/me`, `/resend-verification`, `/logout`, `/logout-all`.
- Background jobs use Asynq (`apps/api/internal/lib/job`). Define new task payloads in `email_tasks.go`, register them in `JobService.Start`, and wire handlers in `handlers.go`.
- Email templates live in `apps/api/templates/emails` and are generated from `packages/emails`.
- OpenAPI docs are written to `apps/api/static/openapi.json` and served at `/api/docs`. Update `packages/zod` and `packages/openapi/src/contracts` when endpoints change.

## Web (apps/web)

- Vite + React + TypeScript with routing in `apps/web/src/router.tsx`.
- Route-based pages live in `apps/web/src/pages`.
- Data layer uses `@ts-rest/react-query` with the axios fetcher in `apps/web/src/api/index.ts`.
- UI uses Tailwind + shadcn/ui; components live in `apps/web/src/components/ui`.
- Auth is cookie-based only. The API client uses `withCredentials: true` and retries once after `/api/v1/auth/refresh`.
- Protected routes use `apps/web/src/auth/require-auth.tsx` (calls `/api/v1/auth/me`).
- Auth routes under `/auth`: `/auth/login`, `/auth/register`, `/auth/verify-email`, `/auth/forgot-password`, `/auth/me`.
- Google login uses `@react-oauth/google` (provider in `apps/web/src/main.tsx`).
- Prefer classes and variables from `apps/web/src/index.css`; avoid arbitrary values unless necessary. User-facing text stays in English.

## Packages (packages/\*)

- `@go-kickstart/zod` (`packages/zod`): source of truth for API request/response schemas (exported from `packages/zod/src/index.ts`).
- `@go-kickstart/openapi` (`packages/openapi`): builds the OpenAPI spec from Zod + ts-rest contracts in `packages/openapi/src/contracts`. Regenerate with `bun run openapi:generate`.
- `@go-kickstart/emails` (`packages/emails`): React Email templates in `packages/emails/src/templates`. Export HTML to `apps/api/templates/emails` via `bun run emails:generate`.

## Testing

- Services: unit tests only, mock repositories.
- Repositories: integration tests with real PostgreSQL (Testcontainers), no SQL mocking.
- Handlers: thin HTTP tests only, mock services.
- Tests live next to code (`foo.go` -> `foo_test.go` / `foo_integration_test.go`).
- Use helpers in `apps/api/internal/testing` (`SetupTestDB`, `WithRollbackTransaction`).
