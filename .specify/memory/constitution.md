<!--
Sync Impact Report
- Version change: 1.0.0 -> 1.0.1
- Modified principles: placeholders -> concrete principles
- Added sections: Architecture & Stack Constraints; Development Workflow & Quality Gates
- Removed sections: None
- Templates requiring updates:
  - ✅ .specify/templates/plan-template.md
  - ✅ .specify/templates/spec-template.md
  - ✅ .specify/templates/tasks-template.md
  - ⚠ .specify/templates/commands/ (directory not found)
-
Follow-up TODOs: None
-->
# Go Kickstart Constitution

## Core Principles

### Contract-First Schemas & OpenAPI
All API request/response schemas MUST be defined in `packages/zod` and exported
from `packages/zod/src/index.ts`. When endpoints change, the corresponding
contracts in `packages/openapi/src/contracts` MUST be updated and
`bun run openapi:generate` MUST be run to regenerate
`apps/api/static/openapi.json`. Manual edits to generated OpenAPI artifacts are
prohibited. Rationale: a single source of truth prevents API and documentation
from drifting.

### Layered Go API Responsibilities
Handlers MUST focus on parsing/validation, delegating to services, and returning
errors; services MUST own business rules; repositories MUST be data access only.
New endpoints MUST use `handler.Handle`, `handler.HandleNoContent`, or
`handler.HandleFile`, and use `utils.ParseUUIDParam` for `:id` parameters.
Services MUST return `errs.ErrorResponse` for expected failures and handlers MUST
let `GlobalErrorHandler` format responses. Rationale: consistent behavior and
error mapping across the API.

### Testing Discipline By Layer
When behavior changes or new logic is added, tests MUST be added at the
appropriate layer: services use unit tests with mocked repositories;
repositories use integration tests with real PostgreSQL (Testcontainers); and
handlers use thin HTTP tests with mocked services. Tests MUST be deterministic,
fast, isolated, and avoid real network calls; do not test external libraries.
Rationale: fast feedback while exercising true decision points.

### Cookie-Only Auth & Client Safety
The frontend MUST use cookie-based authentication only and MUST NOT read, decode,
or store JWTs or use localStorage/sessionStorage for auth. The API client MUST
use `withCredentials: true` and retry once after `/api/v1/auth/refresh`.
Protected routes MUST use `apps/web/src/auth/require-auth.tsx` and call
`/api/v1/auth/me`. Rationale: reduces token leakage and keeps auth centralized.

### Generated Artifacts & Monorepo Workflow
Email templates MUST be authored in `packages/emails` and exported via
`bun run emails:generate` into `apps/api/templates/emails`; generated output MUST
NOT be edited manually. Shared functionality MUST live in `packages/*` and the
repo MUST be operated through Bun/Turborepo workspace commands from the root.
Rationale: keeps generated artifacts and cross-app dependencies consistent.

## Architecture & Stack Constraints

- The stack is Go (API) + Vite/React (web) with Bun workspaces and Turborepo.
- API uses Fiber + GORM with PostgreSQL; web uses React Router, React Query, and
  Tailwind + shadcn/ui.
- Auth uses short-lived JWT access tokens, refresh tokens, and `auth_sessions`.
- Caching MUST use `internal/lib/cache` and only when `Config.Cache.TTL > 0`.
- Background jobs MUST use Asynq; new tasks are registered in the job service.
- Containerization is supported via Dockerfiles and `docker-compose.yml`.

## Development Workflow & Quality Gates

- Database schema changes MUST use the migration commands in `apps/api`.
- API endpoint changes MUST update Zod schemas, OpenAPI contracts, and regenerate
  the OpenAPI spec before merge.
- Email template changes MUST regenerate HTML outputs before merge.
- Request logging MUST include the request ID and handlers SHOULD use the shared
  logging helpers.
- PRs MUST pass constitution checks and include tests per the testing principle
  unless a waiver is documented in the spec with rationale.

## Governance

- This constitution supersedes other development guidance.
- Amendments require a PR that updates this file, includes rationale and
  migration/rollout notes, updates dependent templates/docs, and refreshes the
  Sync Impact Report.
- Versioning follows semantic versioning: MAJOR for principle removals or
  redefinitions, MINOR for new principles/sections or materially expanded
  guidance, PATCH for clarifications or typo fixes.
- Reviewers MUST verify compliance for every PR and document any approved
  waivers in the feature spec. Runtime guidance lives in `AGENTS.md`.

**Version**: 1.0.1 | **Ratified**: 2026-01-31 | **Last Amended**: 2026-01-31
