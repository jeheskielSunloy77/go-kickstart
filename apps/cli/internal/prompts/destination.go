package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func DestinationFlow(cfg *scaffold.ScaffoldConfiguration) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Base destination path").Value(&cfg.Destination).Placeholder("current directory"),
		),
	)
	form.WithTheme(huh.ThemeCharm())
	form.WithWidth(80)
	form.WithHeight(10)
	form.WithOutput(os.Stdout)
	return form.Run()
}
