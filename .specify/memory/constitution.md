<!--
Sync Impact Report
- Version change: N/A -> 1.0.0
- Modified principles: placeholders -> concrete scaffolder principles
- Added sections: Quality Gates; Development Workflow
- Removed sections: None
- Templates requiring updates:
  - ✅ updated: .specify/templates/plan-template.md
  - ✅ updated: .specify/templates/tasks-template.md
  - ✅ no change needed: .specify/templates/spec-template.md
  - ⚠ pending: .specify/templates/commands/ (directory not found in this repo)
- Follow-up TODOs: None
-->
# Go Kickstart CLI Scaffolder Constitution

## Core Principles

### Deterministic, Safe Scaffolding (NON-NEGOTIABLE)
The scaffolder MUST produce consistent output given the same inputs.
It MUST NOT overwrite a non-empty destination directory without explicit user
confirmation (interactive) and MUST fail fast (non-interactive).
Rationale: prevents data loss and makes behavior testable and trustworthy.

### Testing At Decision Points (NON-NEGOTIABLE)
When new behavior is added or existing behavior changes, tests MUST be added.
Tests MUST be deterministic, isolated, and fast:
- Prefer unit tests for validation, config mapping, inclusion/exclusion rules,
  env merging, and rendering transforms.
- Use `t.TempDir()` and `testing/fstest` (or equivalent) to avoid mutating real
  user data.
- Avoid real network calls.
Rationale: scaffolding bugs are costly; decision-point tests catch regressions.

### Consistent UX and Error Clarity
Interactive flows MUST be consistent and predictable:
- Provide clear defaults and a review step before writing files.
- Use consistent wording, spacing, and styling across prompts.
- Errors MUST be actionable (what failed, why, and what the user can do next).
Non-interactive mode MUST fail with clear messages and a non-zero exit code.
Rationale: a scaffolder is a user-facing product; UX consistency reduces churn.

### Performance and Resource Discipline
The default scaffold MUST complete in under 2 minutes on a typical developer
machine.
The scaffolder SHOULD avoid unnecessary work (e.g., rewriting unchanged content,
copying ignored artifacts) and SHOULD keep memory usage reasonable by streaming
file copies where practical.
Rationale: fast scaffolding encourages adoption and enables automation.

### Maintainable, Readable Code
Code MUST be kept simple and maintainable:
- Prefer explicit inputs/outputs over hidden global state.
- Keep packages cohesive (validation, prompting, rendering, templating).
- Run gofmt on Go code; avoid cleverness that complicates future changes.
Rationale: the scaffolder will evolve; clarity reduces bug density.

## Quality Gates

- All PRs MUST keep `go test ./...` passing for the CLI module.
- Behavior changes MUST include corresponding tests unless a waiver is
  documented in the feature spec with rationale.
- Template changes MUST keep the scaffolded project runnable (update the
  template README/AGENTS as needed).

## Development Workflow

- Prefer small PRs that keep generated output stable and reviewable.
- When changing the scaffolded project template, update:
  - `apps/cli/templates/monorepo/README.md`
  - `apps/cli/templates/monorepo/AGENTS.md`
- When changing the CLI contract, update:
  - `specs/*/contracts/` (if present)
  - root `README.md` if user-facing behavior changes

## Governance

This constitution supersedes other development guidance for the scaffolder repo.
Amendments require:
- Updating this file with a Sync Impact Report
- Updating dependent templates/docs when referenced
- Following semantic versioning for constitution version updates:
  - MAJOR: principle removals/redefinitions
  - MINOR: new principles/sections or materially expanded rules
  - PATCH: clarifications and non-semantic refinements
Reviewers MUST verify compliance for each PR and document approved waivers in
the feature spec.

**Version**: 1.0.0 | **Ratified**: 2026-01-31 | **Last Amended**: 2026-01-31
