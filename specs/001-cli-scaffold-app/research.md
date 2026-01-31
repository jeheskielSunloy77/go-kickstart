# Phase 0 Research: CLI Scaffold App

## Decision 1: CLI Framework
- **Decision**: Use Cobra for command structure and subcommands.
- **Rationale**: Aligns with standard Go CLI patterns, supports flags, help, and nested commands with minimal boilerplate.
- **Alternatives considered**: urfave/cli, kong.

## Decision 2: Interactive TUI Stack
- **Decision**: Use Charmbracelet stack: `huh` for prompts, `lipgloss` for styling, and `bubbles/spinner` for progress.
- **Rationale**: Produces a modern, polished TUI experience consistent with the “beautiful” requirement, with cohesive styling and composable components.
- **Alternatives considered**: `survey/v2` for prompts; `briandowns/spinner` for progress.

## Decision 3: Template Packaging
- **Decision**: Store template files in `apps/cli/templates/monorepo/` and embed them in the CLI binary.
- **Rationale**: Embedding avoids external template dependencies and makes distribution portable while keeping templates co-located in the repo.
- **Alternatives considered**: Loading templates from filesystem at runtime only.

**Template contents**: The scaffold source includes the templatized monorepo layout (e.g., `apps/api`, `apps/web`, `packages/*`, `docker-compose.yml`, and other repo root files) under `apps/cli/templates/monorepo/`.

## Decision 4: Command Modes
- **Decision**: Provide both `gokickstart new <name> [path] [flags]` (non-interactive) and `gokickstart new` / `--interactive` (wizard).
- **Rationale**: Satisfies automation workflows and the default interactive experience.
- **Alternatives considered**: Interactive-only; configuration-file-only.

## Decision 5: Defaults and Review
- **Decision**: Always show a review screen with defaults and allow edits before generation.
- **Rationale**: Prevents mistakes and keeps the TUI flow predictable and friendly.
- **Alternatives considered**: No review step; advanced-mode toggle.

## Decision 6: Git Initialization
- **Decision**: If Git is selected, initialize a repo and create an initial commit.
- **Rationale**: Creates a clean baseline for users immediately after scaffold.
- **Alternatives considered**: Initialize without commit; leave Git setup to user.
