package prompts

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/validate"
)

type ReviewAction string

const (
	ReviewGenerate ReviewAction = "generate"
	ReviewEdit     ReviewAction = "edit"
	ReviewCancel   ReviewAction = "cancel"
)

func ReviewConfig(cfg scaffold.ScaffoldConfiguration) (ReviewAction, error) {
	summary := reviewSummary(cfg)

	action := ReviewGenerate
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(ui.ReviewTitle).Description(summary),
			huh.NewSelect[ReviewAction]().Title(ui.ReviewActionTitle).
				Options(
					huh.NewOption(ui.ReviewGenerateLabel, ReviewGenerate),
					huh.NewOption(ui.ReviewEditLabel, ReviewEdit),
					huh.NewOption(ui.ReviewCancelLabel, ReviewCancel),
				).Value(&action),
		),
	)
	form.WithTheme(ui.HuhTheme())
	form.WithWidth(80)
	form.WithHeight(20)
	form.WithOutput(os.Stdout)

	return action, form.Run()
}

func resolveDisplayDestination(cfg scaffold.ScaffoldConfiguration) string {
	if resolved, err := validate.ResolveProjectDestination(cfg.Destination, cfg.ProjectName); err == nil {
		return resolved
	}
	return cfg.Destination
}

func reviewSummary(cfg scaffold.ScaffoldConfiguration) string {
	return fmt.Sprintf(
		"%s\nProject name: %s\nDestination: %s\nModule: %s\nWeb: %t\nDocker: %t\nGit: %t\nDatabase: %s\nStorage: %s\n",
		ui.ReviewSummaryHeading,
		cfg.ProjectName,
		resolveDisplayDestination(cfg),
		cfg.ModulePath,
		cfg.IncludeWeb,
		cfg.IncludeDocker,
		cfg.InitGit,
		cfg.DatabaseType,
		cfg.Storage.Type,
	)
}
