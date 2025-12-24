# Project Context

- Monorepo managed by Turborepo + Bun workspaces.
- Go backend lives in `apps/api`; shared TypeScript packages live in `packages/*`.
- `apps/web` exists (Next.js + Tailwind + shadcn + zod) but is intentionally omitted here.
- OpenAPI docs are generated from shared Zod schemas in `packages/zod` and written into `apps/api/static/openapi.json`.
- Email templates are authored in React Email (`packages/emails`) and exported to Go HTML templates consumed by the API.

## App #1: API (apps/api)

### Architecture & Flow

- Clean layers: handlers -> services -> repositories -> models.
- Entry points: `apps/api/cmd/api/main.go` (server), `apps/api/cmd/seed/main.go` (seeder).
- Routes are registered in `apps/api/internal/router/routes.go`; middleware order in `apps/api/internal/router/router.go`.
- Prefer `handler.Handle` / `handler.HandleNoContent` / `handler.HandleFile` for new endpoints (validation, logging, tracing).

### Requests, Validation, Errors

- Request DTOs implement `validation.Validatable` using go-playground/validator tags.
- Use `validation.BindAndValidate` or the `handler.Handle` wrappers (they call it for you).
- Use `utils.ParseUUIDParam` for `:id` params to standardize 400s.
- Services return `errs.HTTPError` for expected failures; wrap DB errors with `sqlerr.HandleError`.
- Handlers should return errors and let `GlobalErrorHandler` format responses (avoid manual error JSON in handlers).

### Data & Migrations

- PostgreSQL + GORM; migrations live in `apps/api/internal/database/migrations`.
- Use `make migrations-new NAME=...` / `make migrations-up` / `make migrations-down` in `apps/api`.
- Repositories are data access only; services implement business rules and validations.
- Use the generic `ResourceRepository`/`ResourceService`/`ResourceHandler` when a model fits CRUD patterns.

### Auth, Context, and Logging

- JWT auth via `middleware.Auth.RequireAuth`; sets `user_id` in Fiber locals.
- Request ID is set in middleware and injected into logs; use `middleware.GetLogger` in handlers.
- Context timeouts should use `server.Config.Server.ReadTimeout` / `WriteTimeout`.

### Jobs & Emails

- Background jobs use Asynq (`apps/api/internal/lib/job`).
- New task types: define payload + task in `email_tasks.go`, register in `JobService.Start`, and wire handlers in `handlers.go`.
- Email sending uses Go templates in `apps/api/templates/emails` via `internal/lib/email`.

### OpenAPI

- `apps/api/static/openapi.json` and `/api/docs` are generated from `packages/openapi`.
- Update `packages/zod` and `packages/openapi/src/contracts` when adding or changing endpoints.
- Regenerate with `bun run openapi:generate` at repo root.

### Testing Guidelines

#### Goals

- Test behavior where decisions are made.
- Prefer fewer, high-value tests over many shallow tests.
- Keep tests deterministic, fast, and free of global state.
- Avoid real network calls.
- No hidden magic; keep setup explicit.

#### What To Test

- Services:
  - Unit tests only.
  - Mock repositories.
  - Cover business rules, validation, and error handling.
- Repositories:
  - Integration tests only.
  - Use a real PostgreSQL database (Testcontainers).
  - No SQL mocking.
- Handlers:
  - Thin HTTP tests only.
  - Mock services.
  - Test request parsing, status codes, and error mapping.

#### What Not To Test

- Do not unit-test repositories with mocks.
- Do not test business logic in handlers.
- Do not duplicate service tests at handler level.
- Do not test Fiber or any other external libraries functionality.

#### Style & Structure

- Prefer table-driven tests.
- Tests live next to code: `foo.go` -> `foo_test.go` or `foo_integration_test.go`.
- Use helpers in `apps/api/internal/testing` (`SetupTestDB`, `WithRollbackTransaction`) for integration tests.

## Packages (packages/\*)

### @go-kickstart/zod (packages/zod)

- Source of truth for API request/response schemas.
- Update when API models or validation rules change.
- Exported from `packages/zod/src/index.ts`.

### @go-kickstart/openapi (packages/openapi)

- Builds the OpenAPI spec from Zod + ts-rest contracts. use `bun run openapi:generate` to regenerate.
- Contracts live in `packages/openapi/src/contracts`; use `createResourceContract` for CRUD resources.
- Everytime a route is added/changed in the API, update the corresponding contract here.

### @go-kickstart/emails (packages/emails)

- React Email templates live in `packages/emails/src/templates`.
- Use Go template placeholders (e.g., `{{.UserFirstName}}`) to match `internal/lib/email` data keys.
- Export HTML to `apps/api/templates/emails` via `bun run emails:generate`.
