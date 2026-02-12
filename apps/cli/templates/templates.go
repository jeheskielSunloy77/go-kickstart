package templates

import "embed"

//go:embed monorepo/** monorepo/apps/web/.env monorepo/apps/web/.env.example monorepo/apps/api/.env.example.tmpl
var MonorepoFS embed.FS
