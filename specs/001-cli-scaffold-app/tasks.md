# Tasks: CLI Scaffold App

**Input**: Design documents from `/specs/001-cli-scaffold-app/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are REQUIRED when behavior changes or new logic is added, per
the constitution.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the CLI app skeleton and lay out template location.

- [x] T001 Create CLI module skeleton at `apps/cli/` (`apps/cli/go.mod`, `apps/cli/cmd/gokickstart/main.go`)
- [x] T002 Add Cobra root command and wire main entrypoint in `apps/cli/internal/cmd/root.go`
- [x] T003 [P] Add TUI UX foundation (styles, theme tokens) in `apps/cli/internal/ui/styles.go`
- [x] T004 [P] Add dependencies (cobra + charmbracelet stack) and run `go mod tidy` in `apps/cli/go.mod`
- [x] T005 Create template source directory `apps/cli/templates/monorepo/` with placeholder `.keep` file

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core scaffolding engine + validation + template packaging. MUST complete before any user story.

- [x] T006 Move current scaffold source into templates: create `apps/cli/templates/monorepo/apps/` and move `apps/api/` and `apps/web/` into it
- [x] T007 Move current scaffold source into templates: move `packages/` into `apps/cli/templates/monorepo/packages/`
- [x] T008 Move current scaffold source into templates: move `docker-compose.yml` into `apps/cli/templates/monorepo/docker-compose.yml`
- [x] T009 Define template include/exclude rules (ignore `.git/`, `node_modules/`, build artifacts) in `apps/cli/internal/scaffold/ignore.go`
- [x] T010 Embed templates using `go:embed` in `apps/cli/templates/templates.go` (root: `apps/cli/templates/monorepo/**`)
- [x] T011 Define `ScaffoldConfiguration` and option enums in `apps/cli/internal/scaffold/config.go`
- [x] T012 [P] Implement validation for project name and module path in `apps/cli/internal/validate/inputs.go`
- [x] T013 Implement destination path resolution (2nd arg defaulting to `.`) and safety checks in `apps/cli/internal/validate/path.go`
- [x] T014 Implement “do not overwrite non-empty directories without confirmation” helper in `apps/cli/internal/scaffold/safety.go`
- [x] T015 Implement filesystem renderer that can render from an `fs.FS` into an OS dir in `apps/cli/internal/scaffold/renderer.go`
- [x] T016 Implement conditional inclusion (web/docker/storage selection) in `apps/cli/internal/scaffold/conditions.go`
- [x] T017 Implement text templating pass for targeted files (Go module path, project name, env values) in `apps/cli/internal/scaffold/templating.go`
- [x] T018 Implement `.env` generation from `.env.example` + overrides (preserve comments as best-effort) in `apps/cli/internal/scaffold/env.go`
- [x] T019 Implement Git init + initial commit helper in `apps/cli/internal/scaffold/git.go`
- [x] T020 [P] Add unit tests for validation helpers in `apps/cli/internal/validate/inputs_test.go`
- [x] T021 [P] Add unit tests for destination-path safety behavior in `apps/cli/internal/validate/path_test.go`
- [x] T022 Add scaffold engine tests using a small in-memory fixture FS in `apps/cli/internal/scaffold/renderer_test.go`

**Checkpoint**: Foundation ready — user story work can begin.

---

## Phase 3: User Story 1 - Create a New Project from CLI (Priority: P1) 🎯 MVP

**Goal**: Interactive `gokickstart new` wizard that scaffolds a project with defaults (basic flow).

**Independent Test**: Run `go run ./apps/cli/cmd/gokickstart new` and complete the wizard by entering name + path + module; verify output directory exists and includes expected default components.

### Tests for User Story 1

- [x] T023 [P] [US1] Add CLI integration test that runs scaffold against a temp dir using a fixture template FS in `apps/cli/internal/cmd/new_basic_test.go`

### Implementation for User Story 1

- [x] T024 [US1] Implement `new` Cobra command shell and args contract in `apps/cli/internal/cmd/new.go`
- [x] T025 [US1] Implement interactive basic flow prompts (name, destination path, module) in `apps/cli/internal/prompts/basic.go`
- [x] T026 [US1] Implement “use defaults” fast-path and map to config defaults in `apps/cli/internal/scaffold/defaults.go`
- [x] T027 [US1] Add review screen that displays defaults and allows edits before generation in `apps/cli/internal/prompts/review.go`
- [x] T028 [US1] Wire spinner/progress output during file generation in `apps/cli/internal/ui/progress.go`
- [x] T029 [US1] Print success summary + next steps in `apps/cli/internal/ui/summary.go`

**Checkpoint**: US1 works independently (interactive basic scaffold).

---

## Phase 4: User Story 2 - Choose What Gets Included (Priority: P2)

**Goal**: Interactive advanced flow for selecting components and editing most variables (web/docker/git/storage/db details).

**Independent Test**: Run `go run ./apps/cli/cmd/gokickstart new`, choose advanced flow, toggle web/docker/storage, enter DB and (optional/required) storage details, confirm review screen, and verify output matches selections.

### Tests for User Story 2

- [x] T030 [P] [US2] Add scaffold engine tests for conditional inclusion (web/docker/storage) using fixture templates in `apps/cli/internal/scaffold/conditions_test.go`
- [x] T031 [P] [US2] Add `.env` generation tests for postgres + local storage defaults and overrides in `apps/cli/internal/scaffold/env_test.go`

### Implementation for User Story 2

- [x] T032 [US2] Add advanced-flow prompt entry point and navigation (Back/Next) in `apps/cli/internal/prompts/flow.go`
- [x] T033 [US2] Implement component selection prompts (web, docker, git) in `apps/cli/internal/prompts/components.go`
- [x] T034 [US2] Implement DB selection and postgres connection details prompts (defaults prefilled) in `apps/cli/internal/prompts/db.go`
- [x] T035 [US2] Implement storage selection (local vs s3) and required/optional details rules in `apps/cli/internal/prompts/storage.go`
- [x] T036 [US2] Ensure review screen can edit advanced settings without losing prior answers in `apps/cli/internal/prompts/review.go`
- [x] T037 [US2] Apply selected options to scaffold engine (conditions + env generation + git) in `apps/cli/internal/scaffold/scaffold.go`

**Checkpoint**: US2 works independently (interactive advanced customization).

---

## Phase 5: User Story 3 - Run in Automation-Friendly Mode (Priority: P3)

**Goal**: Non-interactive `gokickstart new <name> [path] [flags]` that fully scaffolds without prompts.

**Independent Test**: Run `go run ./apps/cli/cmd/gokickstart new myapp ./out --module github.com/acme/myapp --no-web --docker --git --storage local --db postgres ...` and confirm it exits 0 and generates the expected structure.

### Tests for User Story 3

- [x] T038 [P] [US3] Add CLI flag parsing + required-flag validation tests in `apps/cli/internal/cmd/new_flags_test.go`

### Implementation for User Story 3

- [x] T039 [US3] Implement positional args: `<name> [path]` and default path to current directory in `apps/cli/internal/cmd/new.go`
- [x] T040 [US3] Implement flags-only non-interactive mode that errors if required values are missing in `apps/cli/internal/cmd/new.go`
- [x] T041 [US3] Ensure `--interactive` forces wizard even if args/flags provided in `apps/cli/internal/cmd/new.go`

**Checkpoint**: US3 works independently (automation-friendly, no prompts).

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Hardening, documentation alignment, and developer ergonomics.

- [x] T042 [P] Align docs to current contract: fix command synopsis in `/home/jay/Programming/go-kickstart/specs/001-cli-scaffold-app/research.md` (Decision 4) to include `[path]`
- [x] T043 [P] Add `--version` and `version` output in `apps/cli/internal/cmd/root.go`
- [x] T044 Improve error messages (actionable, consistent) in `apps/cli/internal/ui/errors.go`
- [x] T045 Add a top-level README section for the CLI and template location in `README.md`
- [x] T046 Add a quickstart validation script/task note and run it manually: `/home/jay/Programming/go-kickstart/specs/001-cli-scaffold-app/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Phase 1
- **User Stories (Phases 3-5)**: Depend on Phase 2
- **Polish (Phase 6)**: Depends on whichever user stories you intend to ship

### User Story Dependencies

- **US1 (P1)**: Depends on Phase 2 only
- **US2 (P2)**: Depends on Phase 2; builds on the same scaffold engine as US1
- **US3 (P3)**: Depends on Phase 2; uses the same scaffold engine as US1

### Parallel Opportunities

- Phase 1: T003 and T004 can run in parallel.
- Phase 2: T012, T020, T021, and T022 can be parallelized after T011/T013 exist.
- US2 tests (T030, T031) can run in parallel.

---

## Parallel Example: User Story 2

```bash
# Run these in parallel (different files):
Task: "Add scaffold engine tests for conditional inclusion (web/docker/storage) using fixture templates in apps/cli/internal/scaffold/conditions_test.go"
Task: "Add .env generation tests for postgres + local storage defaults and overrides in apps/cli/internal/scaffold/env_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 and Phase 2
2. Complete Phase 3 (US1)
3. Validate via the US1 independent test

### Incremental Delivery

1. Add US2 advanced flow
2. Add US3 non-interactive mode
3. Finish polish tasks
