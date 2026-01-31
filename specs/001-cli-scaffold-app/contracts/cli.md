# CLI Contract: gokickstart

## Commands

### `gokickstart new <name> [path] [flags]`

Non-interactive mode. All required inputs are provided via arguments and flags. `path` defaults to the current directory when omitted. If any required input is missing, the command fails with a clear error.

### `gokickstart new`

Interactive mode. Launches a step-by-step TUI wizard with Back/Next navigation and a final review screen. The wizard asks for the destination path and defaults it to the current directory. `--interactive` also forces this mode.
it have 2 flows, first one is the basic flow where the user only provide the project name, destination path and the module path, and all other options are set to defaults. The second flow is the advanced flow where the user can choose which parts of the template are included and adjust most template variables.

## Arguments and Flags (non-interactive)

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

## Exit Behavior

- Success: prints a summary and next steps.
- Failure: prints a clear error and exits non-zero.
