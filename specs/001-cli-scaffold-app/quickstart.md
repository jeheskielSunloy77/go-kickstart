# Quickstart: CLI Scaffold App

## Run Locally (Interactive)

1. From repo root, run the CLI entry point:
   - `go run ./apps/cli/cmd/gokickstart new`
2. Follow the TUI wizard steps and confirm on the review screen.
   - The wizard asks for the destination path and defaults it to the current directory.
3. Review the success summary and next steps.

## Run Locally (Non-Interactive)

1. From repo root, run:
   - `go run ./apps/cli/cmd/gokickstart new <name> [path] [flags]`
2. Provide all required arguments/flags (name, module, db, storage, etc.). `path` defaults to the current directory when omitted.

## Output

- A new project folder is created at the destination path with the selected components.
- Environment files are generated from defaults and user overrides.
- Optional Git initialization includes an initial commit.

## Quickstart Validation

- Run the interactive flow and confirm a project is created in a temp directory.
- Run the non-interactive flow with flags and confirm it exits 0 and generates the expected files.
