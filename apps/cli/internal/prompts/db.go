package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/ui"
)

func DatabaseFlow(cfg *scaffold.ScaffoldConfiguration) error {
	cfg.DatabaseType = scaffold.DatabasePostgres
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title(ui.DBHostTitle).Description(ui.DBHostDesc).Value(&cfg.DBConnection.Host),
			huh.NewInput().Title(ui.DBPortTitle).Description(ui.DBPortDesc).Value(&cfg.DBConnection.Port),
			huh.NewInput().Title(ui.DBUserTitle).Description(ui.DBUserDesc).Value(&cfg.DBConnection.User),
			huh.NewInput().Title(ui.DBPassTitle).Description(ui.DBPassDesc).Value(&cfg.DBConnection.Password),
			huh.NewInput().Title(ui.DBNameTitle).Description(ui.DBNameDesc).Value(&cfg.DBConnection.Name),
			huh.NewInput().Title(ui.DBSSLTitle).Description(ui.DBSSLDesc).Value(&cfg.DBConnection.SSLMode),
		),
	)
	form.WithTheme(ui.HuhTheme())
	form.WithWidth(80)
	form.WithHeight(20)
	form.WithOutput(os.Stdout)
	return form.Run()
}
