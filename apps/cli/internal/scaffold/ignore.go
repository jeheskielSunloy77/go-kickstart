package scaffold

var DefaultIgnoreGlobs = []string{
	"**/.git/**",
	"**/node_modules/**",
	"**/.turbo/**",
	"**/dist/**",
	"**/build/**",
	"**/out/**",
	"**/*.log",
	"**/.env",
	"**/.env.*",
	"**/*.tsbuildinfo",
}
