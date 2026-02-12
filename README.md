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

## CLI Contract

### Commands

#### `gokickstart new <name> [path] [flags]`

Non-interactive mode. All required inputs are provided via arguments and flags. `path` defaults to the current directory when omitted. If any required input is missing, the command fails with a clear error.

#### `gokickstart new`

Interactive mode. Launches a step-by-step TUI wizard with Back/Next navigation and a final review screen. The wizard asks for the destination path and defaults it to the current directory. `--interactive` also forces this mode.
it have 2 flows, first one is the basic flow where the user only provide the project name, destination path and the module path, and all other options are set to defaults. The second flow is the advanced flow where the user can choose which parts of the template are included and adjust most template variables.

### Arguments and Flags (non-interactive)

- `--name` (string): Folder/app name.
- `path` (arg): Destination path for the generated project (defaults to current directory).
- `--module` (string): Go module path (e.g., `github.com/acme/foo`).
- `--web` / `--no-web`: Include or exclude `apps/web`.
- `--db` (enum): `postgres` (only option in this release).
- `--db-host`, `--db-port`, `--db-user`, `--db-password`, `--db-name`, `--db-ssl-mode`: Database connection details (required when `--db=postgres`).
- `--pkg` (enum): `bun` (only option in this release).
- `--docker` / `--no-docker`: Include or exclude Docker Compose.
- `--git` / `--no-git`: Initialize Git repo and create an initial commit.
- `--storage` (enum): `local` or `s3`.
- `--s3-endpoint`, `--s3-region`, `--s3-bucket`, `--s3-access-key`, `--s3-secret-key`: Required when `--storage=s3`.
- `--interactive`: Force interactive wizard.

### Exit Behavior

- Success: prints a summary and next steps.
- Failure: prints a clear error and exits non-zero.

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
