package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func ComponentsFlow(cfg *scaffold.ScaffoldConfiguration) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Include web app").Value(&cfg.IncludeWeb),
			huh.NewConfirm().Title("Include Docker Compose").Value(&cfg.IncludeDocker),
			huh.NewConfirm().Title("Initialize git repository").Value(&cfg.InitGit),
		),
	)
	form.WithTheme(huh.ThemeCharm())
	form.WithWidth(80)
	form.WithHeight(15)
	form.WithOutput(os.Stdout)
	return form.Run()
}
