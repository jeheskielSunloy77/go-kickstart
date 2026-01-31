# Implementation Plan: CLI Scaffold App

**Branch**: `001-cli-scaffold-app` | **Date**: January 31, 2026 | **Spec**: /home/jay/Programming/go-kickstart/specs/001-cli-scaffold-app/spec.md  
**Input**: Feature specification from `/specs/001-cli-scaffold-app/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a Go-based CLI that scaffolds this monorepo using embedded templates in `apps/cli/templates/monorepo/`, providing a modern TUI wizard (interactive by default) and a flags-only non-interactive mode. Use Cobra for command structure and Charmbracelet components for prompts, styling, and progress feedback, while producing deterministic output and honoring customization inputs (name, module, components, storage, db, etc.).

## Technical Context

**Language/Version**: Go 1.22+  
**Primary Dependencies**: Cobra CLI, Charmbracelet `huh` (prompts), `lipgloss` (styling), `bubbles/spinner` (progress)  
**Storage**: Embedded template files and filesystem output  
**Testing**: Go `testing` package with deterministic, filesystem-isolated tests  
**Target Platform**: macOS, Linux, Windows developer machines  
**Project Type**: CLI tool within the monorepo  
**Performance Goals**: Default scaffold completes in under 2 minutes  
**Constraints**: No network required; flags-only non-interactive mode; do not overwrite non-empty directories without confirmation; templates live in `apps/cli/templates/monorepo/` and are embedded into the binary  
**Scale/Scope**: Single project scaffold per invocation

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Contract-first schemas: No API changes planned; N/A for this feature.
- Layered API: No API changes planned; N/A for this feature.
- Testing discipline: CLI logic will be covered with focused unit tests; no layer conflicts.
- Web auth safety: No auth changes; N/A.
- Generated artifacts: No email or OpenAPI generation changes; N/A.
- Monorepo workflow: CLI will live in the monorepo with Bun/Turbo conventions respected.

## Project Structure

### Documentation (this feature)

```text
specs/001-cli-scaffold-app/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
apps/
└── cli/
    ├── cmd/
    │   └── gokickstart/
    ├── internal/
    │   ├── prompts/
    │   ├── render/
    │   ├── scaffold/
    │   └── validate/
    └── templates/
        └── monorepo/
            ├── apps/
            │   ├── api/
            │   └── web/
            ├── packages/
            ├── docker-compose.yml
            └── [other templatized repo files]
```

**Structure Decision**: The repo becomes a CLI-only project inside the existing monorepo structure, with the former `apps/api`, `apps/web`, and `packages/*` moved into `apps/cli/templates/monorepo/` as the scaffold source. The CLI entry point is `apps/cli/cmd/gokickstart`, shared logic lives in `apps/cli/internal/*`, and templates are embedded for portability.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations.

## Constitution Check (Post-Design)

- Contract-first schemas: No API changes planned; still N/A.
- Layered API: No API changes planned; still N/A.
- Testing discipline: Plan includes CLI-focused tests; no conflicts.
- Web auth safety: No auth changes; still N/A.
- Generated artifacts: No changes to OpenAPI or emails; still N/A.
- Monorepo workflow: CLI and templates remain within repo conventions.
