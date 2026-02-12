package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func defaultModulePath(projectName string) string {
	return "github.com/yourorg/" + projectName
}

func BasicFlow(defaults scaffold.ScaffoldConfiguration) (scaffold.ScaffoldConfiguration, error) {
	cfg := defaults
	if cfg.ProjectName == "" {
		cfg.ProjectName = "my-app"
	}
	initialProjectName := cfg.ProjectName

	if cfg.ModulePath == "" {
		cfg.ModulePath = defaultModulePath(cfg.ProjectName)
	}
	initialModulePath := cfg.ModulePath
	shouldSyncDefaultModule := initialModulePath == defaultModulePath(initialProjectName)

	projectForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Project name").Value(&cfg.ProjectName),
		),
	)
	projectForm.WithTheme(huh.ThemeCharm())
	projectForm.WithWidth(80)
	projectForm.WithHeight(10)
	projectForm.WithOutput(os.Stdout)

	if err := projectForm.Run(); err != nil {
		return cfg, err
	}

	if shouldSyncDefaultModule && cfg.ModulePath == initialModulePath {
		cfg.ModulePath = defaultModulePath(cfg.ProjectName)
	}

	moduleForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Go module path").Value(&cfg.ModulePath),
		),
	)
	moduleForm.WithTheme(huh.ThemeCharm())
	moduleForm.WithWidth(80)
	moduleForm.WithHeight(10)
	moduleForm.WithOutput(os.Stdout)

	if err := moduleForm.Run(); err != nil {
		return cfg, err
	}

	return cfg, nil
}
