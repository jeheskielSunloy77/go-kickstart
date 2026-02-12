package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
)

type FlowChoice string

const (
	FlowBasic    FlowChoice = "basic"
	FlowAdvanced FlowChoice = "advanced"
)

func ChooseFlow() (FlowChoice, error) {
	choice := FlowBasic
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[FlowChoice]().
				Title(ui.FlowTitle).
				Description(ui.FlowDescription).
				Options(
					huh.NewOption(ui.FlowBasicLabel, FlowBasic),
					huh.NewOption(ui.FlowAdvLabel, FlowAdvanced),
				).Value(&choice),
		),
	)
	form.WithTheme(ui.HuhTheme())
	form.WithWidth(80)
	form.WithHeight(12)
	form.WithOutput(os.Stdout)
	return choice, form.Run()
}
