package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
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
			huh.NewSelect[FlowChoice]().Title("Choose setup flow").
				Options(
					huh.NewOption("Basic (use defaults)", FlowBasic),
					huh.NewOption("Advanced (customize)", FlowAdvanced),
				).Value(&choice),
		),
	)
	form.WithTheme(huh.ThemeCharm())
	form.WithWidth(80)
	form.WithHeight(12)
	form.WithOutput(os.Stdout)
	return choice, form.Run()
}
