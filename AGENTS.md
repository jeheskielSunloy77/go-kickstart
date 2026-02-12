# Project Context (CLI Scaffolder)

This repository is the Go Kickstart **CLI scaffolder**. Its job is to generate a full monorepo (API + web + shared TS packages) from an embedded template.

The generated project's documentation and engineering rules live in the template:

- `apps/cli/templates/monorepo/README.md` (README for scaffolded projects)
- `apps/cli/templates/monorepo/AGENTS.md` (agent context for scaffolded projects)

## Repo Layout

```
go-kickstart/
├── apps/cli/
│   ├── cmd/gokickstart/          # CLI entrypoint
│   ├── internal/                 # CLI implementation
│   └── templates/monorepo/       # Scaffold source (templatized monorepo)
└── specs/                        # Specs/plans/tasks for the scaffolder itself
```

## What To Edit (Rule of Thumb)

- If you are changing **what users get when they scaffold** (API/web/packages/Docker/etc), edit:
  - `apps/cli/templates/monorepo/**`
  - and keep `apps/cli/templates/monorepo/README.md` + `apps/cli/templates/monorepo/AGENTS.md` in sync.

- If you are changing **how scaffolding works** (prompts, flags, rendering, validation, env generation, git init), edit:
  - `apps/cli/internal/**`
  - `apps/cli/cmd/gokickstart/main.go`

## CLI Contract

- Interactive mode (default): `gokickstart new`
  - Modern TUI wizard
  - Asks for base destination path (defaults to current directory)
  - Generates at `<base>/<project-name>`
  - Supports basic (defaults) and advanced (customize) flows
- Non-interactive mode: `gokickstart new <name> [path] [flags]`
  - Flags-only (no prompts)
  - `path` is a base directory; final destination is `<path>/<name>`
  - omitted `path` defaults base to current directory

Canonical contract doc: `specs/001-cli-scaffold-app/contracts/cli.md`

## Template Packaging

- Templates are embedded from `apps/cli/templates/monorepo/**` (see `apps/cli/templates/templates.go`).
- Rendering reads from an `fs.FS` and writes to disk with optional transforms.
- `.env` files are generated from `.env.example` files plus overrides derived from user inputs.

## Testing Guidelines (Scaffolder)

Goals:
- Test behavior where decisions are made (validation, inclusion/exclusion rules, env merging, flag parsing).
- Keep tests deterministic and filesystem-isolated (use `t.TempDir()` and `testing/fstest`).
- Avoid real network calls.

Where tests live:
- Go unit tests alongside code in `apps/cli/internal/**` (e.g., `*_test.go`).

## Common Commands

Run from `apps/cli/`:

- `go test ./...`
- `go run ./cmd/gokickstart new`
- `go run ./cmd/gokickstart new <name> [path] --module github.com/acme/<name> --storage local`

## Generated Monorepo Context

When you need details about the scaffolded monorepo (API/web/packages conventions, testing by layer, OpenAPI/email generation, auth rules, etc.), refer to:

- `apps/cli/templates/monorepo/AGENTS.md`
