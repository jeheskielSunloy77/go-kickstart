# Data Model: CLI Scaffold App

## Entities

### ScaffoldConfiguration
Represents the full set of inputs for a scaffold run.

**Fields**
- `projectName` (string, required): Folder/app name.
- `destinationPath` (string, required): Output path.
- `modulePath` (string, required): Go module path.
- `includeWeb` (bool, required): Include `apps/web`.
- `databaseType` (enum, required): `postgres` (only option in this release).
- `dbConnection` (object, required when databaseType=postgres): Connection details to populate environment defaults.
  - `host`, `port`, `user`, `password`, `database`, `sslMode` (strings/number, required but pre-filled with defaults).
- `packageManager` (enum, required): `bun` (only option in this release).
- `includeDocker` (bool, required): Include Docker Compose files.
- `initGit` (bool, required): Initialize Git repository.
- `storageType` (enum, required): `local` or `s3`.
- `storageConfig` (object):
  - For `local`: optional overrides with defaults pre-filled.
  - For `s3`: required connection details (endpoint, region, bucket, access key, secret key).
- `useDefaults` (bool, required): Indicates whether user accepted default options at the start.

**Validation Rules**
- `projectName` must be non-empty and filesystem-safe.
- `modulePath` must match Go module path format.
- `destinationPath` must not be a non-empty directory unless explicitly confirmed.
- `storageConfig` is required and complete when `storageType` is `s3`.

### UserInputSet
Represents the raw answers captured during an interactive or flags-only run.

**Fields**
- `mode` (enum): `interactive` or `nonInteractive`.
- `answers` (map): Key-value set of prompts or flags supplied.

### GeneratedProject
Represents the output of a completed scaffold run.

**Fields**
- `path` (string): Output directory.
- `components` (set): Components included (web, docker, storage type, etc.).
- `envFiles` (set): Environment files generated with defaults and overrides.
- `gitInitialized` (bool): Whether a repo and initial commit exist.

## Lifecycle / State Transitions

- `Draft` → `Reviewed` → `Generated`
  - Draft: inputs captured
  - Reviewed: user confirmed review screen
  - Generated: files written and post-steps completed
