package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
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
			huh.NewInput().
				Title(ui.ProjectNameTitle).
				Description(ui.ProjectNameDescription).
				Value(&cfg.ProjectName),
		),
	)
	projectForm.WithTheme(ui.HuhTheme())
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
			huh.NewInput().
				Title(ui.ModuleTitle).
				Description(ui.ModuleDescription).
				Value(&cfg.ModulePath),
		),
	)
	moduleForm.WithTheme(ui.HuhTheme())
	moduleForm.WithWidth(80)
	moduleForm.WithHeight(10)
	moduleForm.WithOutput(os.Stdout)

	if err := moduleForm.Run(); err != nil {
		return cfg, err
	}

	return cfg, nil
}
