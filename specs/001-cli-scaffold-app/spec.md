# Feature Specification: CLI Scaffold App

**Feature Branch**: `001-cli-scaffold-app`  
**Created**: January 31, 2026  
**Status**: Draft  
**Input**: User description: "currently this project it a repository project template. i want to convert it to a cli app that will scaffold the app customized for the user."

## Clarifications

### Session 2026-01-31

- Q: What depth of customization should the first CLI release support? → A: Extensive: expose most template variables for editing.
- Q: How should users navigate the interactive TUI flow? → A: Wizard: Back/Next with a final review screen.
- Q: How should non-interactive inputs be provided? → A: Flags only (all inputs via CLI flags).
- Q: What should happen when the user opts to initialize Git? → A: Initialize Git and create an initial commit.
- Q: How should defaults be handled for configurable values? → A: Always show a review screen with defaults and allow edits.

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Create a New Project from CLI (Priority: P1)

As a developer, I want to run a command-line tool with a step-by-step TUI wizard so I can start a customized app quickly.

**Why this priority**: This is the core value of the feature and the minimum viable experience.

**Independent Test**: Can be fully tested by running the CLI once and confirming a new project is created with expected defaults.

**Acceptance Scenarios**:

1. **Given** a user runs the CLI in an empty directory, **When** they provide a project name and destination, **Then** a new project folder is created with the base template structure.
2. **Given** a user completes the CLI flow, **When** generation finishes, **Then** they see a clear success message and next steps.

---

### User Story 2 - Choose What Gets Included (Priority: P2)

As a developer, I want to choose which parts of the template are included and adjust most template variables so the generated project matches my needs.

**Why this priority**: Customization is the main reason to use a scaffolder instead of copying a repo.

**Independent Test**: Can be fully tested by selecting a subset of options and verifying only those components exist in the output.

**Acceptance Scenarios**:

1. **Given** the CLI offers customization choices, **When** a user selects a subset of components, **Then** the generated project includes only the selected components.
2. **Given** a user accepts all default options, **When** generation completes, **Then** the output matches the full base template.

---

### User Story 3 - Run in Automation-Friendly Mode (Priority: P3)

As a developer, I want to run the scaffolder non-interactively so I can automate project creation.

**Why this priority**: Automation expands usefulness for teams and repeatable setups.

**Independent Test**: Can be fully tested by running the CLI with predefined inputs and confirming the output matches those inputs.

**Acceptance Scenarios**:

1. **Given** a user provides all required inputs via flags, **When** the CLI runs, **Then** it completes without prompting.

---

### Edge Cases

- What happens when the target directory already exists and is not empty?
- How does the system handle invalid or empty project names?
- How does the system handle invalid module path formats?
- What happens if the user cancels the process mid-run?
- How does the system handle insufficient file permissions in the target location?
- What happens when required connection details are missing for selected services?

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The system MUST provide a command-line entry point to start scaffolding a new project.
- **FR-002**: The system MUST collect a project name and destination path from the user.
- **FR-003**: The system MUST validate the project name and destination path and show actionable errors for invalid input.
- **FR-004**: The system MUST allow users to select which template components are included in the generated project.
- **FR-005**: The system MUST allow users to edit most template variables during scaffolding.
- **FR-006**: The system MUST generate output that reflects the selected components and variable values.
- **FR-007**: The system MUST support a non-interactive mode where all required inputs are provided via CLI flags.
- **FR-008**: The system MUST not modify or overwrite existing non-empty directories without explicit user confirmation.
- **FR-009**: The system MUST provide an interactive TUI wizard with Back/Next navigation and a final review screen before generation.
- **FR-009a**: The system MUST present a review screen that shows defaults and allows edits before generation.
- **FR-010**: The system MUST provide a clear success summary and next-step instructions after generation.
- **FR-011**: The system MUST produce consistent output given the same inputs.
- **FR-012**: The system MUST provide a way to view help and version information.
- **FR-013**: The system MUST accept a destination path as the second CLI argument and default to the current directory when omitted.
- **FR-014**: The system MUST ask for the destination path in interactive mode and default to the current directory.
- **FR-015**: The system MUST store the templatized monorepo under `apps/cli/templates/` for scaffold generation.

### Customization Inputs

- **CI-001**: Users MUST be able to set the project folder/app name.
- **CI-002**: Users MUST be able to set the Go module path.
- **CI-003**: Users MUST be able to include or exclude the web app.
- **CI-004**: Users MUST be able to choose the database option; for this release, only PostgreSQL is supported.
- **CI-005**: When PostgreSQL is selected, the system MUST collect connection details to populate the generated environment configuration with defaults pre-filled and reviewable.
- **CI-006**: Users MUST be able to choose the package manager; for this release, only Bun is supported.
- **CI-007**: Users MUST be able to include or exclude Docker Compose files.
- **CI-008**: Users MUST be able to choose whether to initialize a Git repository, and if selected the system MUST create an initial commit.
- **CI-009**: Users MUST be able to choose a storage option: local or S3-compatible.
- **CI-010**: When S3-compatible storage is selected, the system MUST require connection details and populate the generated environment configuration.
- **CI-011**: When local storage is selected, the system MUST pre-fill defaults and allow review and overrides.

### Key Entities _(include if feature involves data)_

- **Scaffold Configuration**: Captures project name, destination, and selected components.
- **Generated Project**: The created project structure and content derived from the template.
- **User Input Set**: The set of answers or parameters used for a run.

### Assumptions

- The default selection includes all components unless a user opts out.
- The default options are: PostgreSQL for the database, Bun for the package manager, local storage and git initialization enabled. this must used if the user wants to use the default options at the beginning of the flow.
- The CLI targets developers running it locally on their machines.
- The scaffolder does not require user authentication.

## Out of Scope

- Supporting MySQL as a database option in this release.
- Supporting pnpm, npm, or yarn as package manager options in this release.

## Constitution Alignment _(mandatory)_

- **Contract-first schemas**: N/A for this feature.
- **Layered API responsibilities**: N/A for this feature.
- **Testing discipline by layer**: CLI behavior will be covered with focused, deterministic tests aligned to decision points.
- **Cookie-only auth & client safety**: N/A for this feature.
- **Generated artifacts & monorepo workflow**: N/A unless new generation steps are added.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: 90% of first-time users in testing complete project generation without assistance.
- **SC-002**: A default scaffold completes in under 2 minutes on a typical developer machine.
- **SC-003**: 95% of runs with valid inputs produce a usable project without missing files or errors.
- **SC-004**: At least 80% of surveyed users report the customization options meet their needs.
