package prompts

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

type ReviewAction string

const (
	ReviewGenerate ReviewAction = "generate"
	ReviewEdit     ReviewAction = "edit"
	ReviewCancel   ReviewAction = "cancel"
)

func ReviewConfig(cfg scaffold.ScaffoldConfiguration) (ReviewAction, error) {
	summary := fmt.Sprintf(
		"Project name: %s\nDestination: %s\nModule: %s\nWeb: %t\nDocker: %t\nGit: %t\nDatabase: %s\nStorage: %s\n",
		cfg.ProjectName,
		cfg.Destination,
		cfg.ModulePath,
		cfg.IncludeWeb,
		cfg.IncludeDocker,
		cfg.InitGit,
		cfg.DatabaseType,
		cfg.Storage.Type,
	)

	action := ReviewGenerate
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("Review").Description(summary),
			huh.NewSelect[ReviewAction]().Title("What would you like to do?").
				Options(
					huh.NewOption("Generate project", ReviewGenerate),
					huh.NewOption("Edit answers", ReviewEdit),
					huh.NewOption("Cancel", ReviewCancel),
				).Value(&action),
		),
	)
	form.WithTheme(huh.ThemeCharm())
	form.WithWidth(80)
	form.WithHeight(20)
	form.WithOutput(os.Stdout)

	return action, form.Run()
}
