package prompts

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/scaffold"
)

func DatabaseFlow(cfg *scaffold.ScaffoldConfiguration) error {
	cfg.DatabaseType = scaffold.DatabasePostgres
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Postgres host").Value(&cfg.DBConnection.Host),
			huh.NewInput().Title("Postgres port").Value(&cfg.DBConnection.Port),
			huh.NewInput().Title("Postgres user").Value(&cfg.DBConnection.User),
			huh.NewInput().Title("Postgres password").Value(&cfg.DBConnection.Password),
			huh.NewInput().Title("Postgres database").Value(&cfg.DBConnection.Name),
			huh.NewInput().Title("Postgres SSL mode").Value(&cfg.DBConnection.SSLMode),
		),
	)
	form.WithTheme(huh.ThemeCharm())
	form.WithWidth(80)
	form.WithHeight(20)
	form.WithOutput(os.Stdout)
	return form.Run()
}
