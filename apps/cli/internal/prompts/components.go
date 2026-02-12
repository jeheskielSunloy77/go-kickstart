package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func ComponentsFlow(cfg *scaffold.ScaffoldConfiguration) error {
	webForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Include web app").Value(&cfg.IncludeWeb),
		),
	)
	webForm.WithTheme(huh.ThemeCharm())
	webForm.WithWidth(80)
	webForm.WithHeight(10)
	webForm.WithOutput(os.Stdout)
	if err := webForm.Run(); err != nil {
		return err
	}

	dockerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Include Docker Compose").Value(&cfg.IncludeDocker),
		),
	)
	dockerForm.WithTheme(huh.ThemeCharm())
	dockerForm.WithWidth(80)
	dockerForm.WithHeight(10)
	dockerForm.WithOutput(os.Stdout)
	if err := dockerForm.Run(); err != nil {
		return err
	}

	gitForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title("Initialize git repository").Value(&cfg.InitGit),
		),
	)
	gitForm.WithTheme(huh.ThemeCharm())
	gitForm.WithWidth(80)
	gitForm.WithHeight(10)
	gitForm.WithOutput(os.Stdout)
	return gitForm.Run()
}
