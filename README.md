# Go Kickstart (CLI Scaffolder)

Go Kickstart is a CLI that generates a production-ready monorepo (Go API + Vite/React web app + shared TypeScript packages) from an embedded template.

This repository is the scaffolder itself. The README for the generated project lives in the template at `apps/cli/templates/monorepo/README.md`.

## What This Repo Contains

```
go-kickstart/
├── apps/cli/                      # CLI app (this repo's primary code)
│   ├── cmd/gokickstart/           # CLI entrypoint
│   ├── internal/                  # CLI implementation (prompts, scaffold engine, validation)
│   └── templates/monorepo/        # Templatized scaffold source (becomes the generated repo)
└── specs/                         # Feature specs/plans for the scaffolder
```

## Why Use It

- Avoid copying a repo and manually renaming things.
- Get a guided interactive setup (modern TUI wizard) or an automation-friendly CLI command.
- Start from a consistent baseline with batteries included (API, web, jobs, caching, OpenAPI, emails, Docker).

## Usage

The CLI is a Go module under `apps/cli/`. Run it from there:

```bash
cd apps/cli

# Interactive (default)
go run ./cmd/gokickstart new

# Non-interactive
go run ./cmd/gokickstart new <name> [path] --module github.com/acme/<name> --storage local
```

Notes:
- `path` defaults to the current directory when omitted.
- Interactive mode always asks for destination path and defaults it to the current directory.

## Customization (High Level)

Inputs supported by the scaffolder (interactive and/or flags):

- Project name + destination path
- Go module path
- Include/exclude web app (`--web` / `--no-web`)
- Database: PostgreSQL only (for now) + connection details for `.env`
- Package manager: Bun only (for now)
- Include/exclude Docker Compose (`--docker` / `--no-docker`)
- Git init + initial commit (`--git` / `--no-git`)
- Storage: local or S3-compatible + `.env` details (S3 details required when selected)

For the full CLI contract, see `specs/001-cli-scaffold-app/contracts/cli.md`.

## Template (Generated Project)

The scaffold source is the directory `apps/cli/templates/monorepo/`. The CLI embeds these files and writes them to the destination path, applying:

- Conditional inclusion (e.g., exclude `apps/web` when `--no-web`)
- String/token replacements (project name/module)
- `.env` generation from `.env.example` plus user overrides

If you want to understand the generated project in depth, read:

`apps/cli/templates/monorepo/README.md`

## Scaffolded Monorepo Overview (What Users Get)

This is a condensed overview of what the generated project contains.

### Repository Layout

```
go-kickstart/
├── apps/api             # Go API (Fiber)
├── apps/web             # Vite + React frontend
├── packages/zod         # Shared Zod schemas
├── packages/openapi     # OpenAPI generation
├── packages/emails      # React Email templates -> Go HTML templates
└── packages/*           # Other shared packages
```

### Prerequisites (Generated Project)

- Go (API)
- Bun + Node (web + workspaces)
- PostgreSQL (database)
- Redis (jobs/cache)

### Quick Start (Generated Project)

```bash
bun install
cp apps/api/.env.example apps/api/.env
cp apps/web/.env.example apps/web/.env
bun run api:migrate:up
bun dev
```

Docker Compose is also supported by the template when enabled during scaffolding.

### Key Features (Generated Project)

- Go API with clean layering (handlers -> services -> repositories -> models)
- Cookie-based auth with refresh flow
- Background jobs (Asynq) + email templates generation
- OpenAPI generation from shared schemas/contracts
- React web app with React Router, React Query, Tailwind, and shadcn/ui

## Developing This CLI

### Root-level shortcuts

You can run common tasks from the repo root (Bun + Turbo):

```bash
bun run cli:run
bun run build
bun run test
bun run lint
bun run format:check
bun run ci:simulate
```

```bash
cd apps/cli
go test ./...
```

To change what gets generated, edit files under `apps/cli/templates/monorepo/` and re-run the CLI.
