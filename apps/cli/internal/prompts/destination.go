package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
)

func DestinationFlow(cfg *scaffold.ScaffoldConfiguration) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(ui.DestinationTitle).
				Description(ui.DestinationDescription).
				Value(&cfg.Destination).
				Placeholder("current directory"),
		),
	)
	form.WithTheme(ui.HuhTheme())
	form.WithWidth(80)
	form.WithHeight(10)
	form.WithOutput(os.Stdout)
	return form.Run()
}
