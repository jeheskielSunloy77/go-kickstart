package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func BasicFlow(defaults scaffold.ScaffoldConfiguration) (scaffold.ScaffoldConfiguration, error) {
	cfg := defaults
	if cfg.Destination == "" {
		cfg.Destination = "."
	}
	if cfg.ProjectName == "" {
		cfg.ProjectName = "my-app"
	}
	if cfg.ModulePath == "" {
		cfg.ModulePath = "github.com/yourorg/" + cfg.ProjectName
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Project name").Value(&cfg.ProjectName),
			huh.NewInput().Title("Destination path").Value(&cfg.Destination).Placeholder("."),
		),
		huh.NewGroup(
			huh.NewInput().Title("Go module path").Value(&cfg.ModulePath),
		),
	)
	form.WithTheme(huh.ThemeCharm())
	form.WithWidth(80)
	form.WithHeight(20)
	form.WithOutput(os.Stdout)

	return cfg, form.Run()
}
