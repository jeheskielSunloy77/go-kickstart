package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
)

func ObservabilityFlow(cfg *scaffold.ScaffoldConfiguration) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[scaffold.ObservabilityProvider]().
				Title(ui.ObservabilityTitle).
				Description(ui.ObservabilityDescription).
				Options(
					huh.NewOption(ui.ObservabilityNoneLabel, scaffold.ObservabilityNone),
					huh.NewOption(ui.ObservabilityGrafanaOSSLabel, scaffold.ObservabilityGrafanaOSS),
				).
				Value(&cfg.Observability),
		),
	)
	form.WithTheme(ui.HuhTheme())
	form.WithWidth(80)
	form.WithHeight(12)
	form.WithOutput(os.Stdout)

	return form.Run()
}
