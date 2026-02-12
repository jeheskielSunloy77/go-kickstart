package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
)

func ComponentsFlow(cfg *scaffold.ScaffoldConfiguration) error {
	webForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(ui.IncludeWebTitle).
				Description(ui.IncludeWebDescription).
				Value(&cfg.IncludeWeb),
		),
	)
	webForm.WithTheme(ui.HuhTheme())
	webForm.WithWidth(80)
	webForm.WithHeight(10)
	webForm.WithOutput(os.Stdout)
	if err := webForm.Run(); err != nil {
		return err
	}

	dockerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(ui.IncludeDockerTitle).
				Description(ui.IncludeDockerDesc).
				Value(&cfg.IncludeDocker),
		),
	)
	dockerForm.WithTheme(ui.HuhTheme())
	dockerForm.WithWidth(80)
	dockerForm.WithHeight(10)
	dockerForm.WithOutput(os.Stdout)
	if err := dockerForm.Run(); err != nil {
		return err
	}

	gitForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(ui.IncludeGitTitle).
				Description(ui.IncludeGitDescription).
				Value(&cfg.InitGit),
		),
	)
	gitForm.WithTheme(ui.HuhTheme())
	gitForm.WithWidth(80)
	gitForm.WithHeight(10)
	gitForm.WithOutput(os.Stdout)
	return gitForm.Run()
}
